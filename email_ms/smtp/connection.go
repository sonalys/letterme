package smtp

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Connection represents a TCP connection which can be insecure or encrypted.
type Connection interface {
	UpgradeTLS(config *tls.Config) error
	ReadLine() ([]byte, error)
	ReadBytes(delim byte) ([]byte, error)
	ReadEnvelope() (reader io.Reader, err error)
	SetDeadline(d time.Duration) error
	AddBuffer(data ...interface{})
	Flush() error
	RemoteAddr() string
	Close() error
}

// ConnectionAdapter implements a connection.
type ConnectionAdapter struct {
	ctx     context.Context
	conn    net.Conn
	bufout  *bufio.Writer
	bufin   *bufio.Reader
	timeout time.Duration
	TLS     bool
}

// NewConnection instantiates a new connection controller.
func NewConnection(ctx context.Context, c net.Conn, timeout time.Duration, tls *tls.Config) (Connection, error) {
	conn := &ConnectionAdapter{
		ctx:     ctx,
		conn:    c,
		bufout:  bufio.NewWriter(c),
		bufin:   bufio.NewReader(c),
		timeout: timeout,
	}
	if tls != nil {
		// try to upgrade to tls, doesn't matter if it fails.
		if err := conn.UpgradeTLS(tls); err != nil {
			return nil, err
		}
	}
	conn.SetDeadline(timeout)
	return conn, nil
}

// UpgradeTLS upgrades the socket connection with tls encryption.
func (c *ConnectionAdapter) UpgradeTLS(config *tls.Config) error {
	// wrap c.conn in a new TLS server side connection
	conn := tls.Server(c.conn, config)
	conn.SetDeadline(time.Now().Add(c.timeout))

	if err := conn.Handshake(); err != nil {
		return err
	}

	c.TLS = true
	c.conn = conn
	return nil
}

func (c *ConnectionAdapter) resetTimeout() {
	// we don't wan't to keep refreshing this deadline when the ms is closing.
	if c.ctx.Err() == nil {
		c.SetDeadline(c.timeout)
	}
}

// ReadLine uses ReadBytes('\n') to read an entire line, max len 4096 bytes.
func (c *ConnectionAdapter) ReadLine() ([]byte, error) {
	buf, err := c.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf, "\r"), nil
}

// ReadBytes reads the buffer until it reaches the first 'delim' byte.
func (c *ConnectionAdapter) ReadBytes(delim byte) ([]byte, error) {
	buf, err := c.bufin.ReadBytes(delim)
	if err != nil {
		return nil, err
	}
	c.resetTimeout()
	return buf[:len(buf)-1], nil
}

// maxEnvelopeDataSize oh yes, hardcoded value, 25 MB per envelope buffer
const maxEnvelopeDataSize = 25 * MB

// outputPool pre-allocates buffers for parsing envelopes.
var outputPool = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, maxEnvelopeDataSize)
		return &buf
	},
}

// ReadEnvelope parses the upcoming data buffer.
func (c *ConnectionAdapter) ReadEnvelope() (reader io.Reader, err error) {
	// We allocate 1 buffer of 25 MB for the whole envelope
	buffer := *outputPool.Get().(*[]byte)

	var size uint32
	var endOffset uint32

	for {
		endOffset = size + MB
		if int(size+MB) > len(buffer) {
			endOffset = maxEnvelopeDataSize
		}

		bytesRead, err := c.conn.Read(buffer[size:endOffset])
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, err
			}
			c.Close()
			return nil, err
		}
		size = size + uint32(bytesRead)

		if bytes.Equal(buffer[size-3:size-1], []byte{'.', '\r'}) {
			return bytes.NewBuffer(buffer[:size]), nil
		}

		if bytesRead == 0 {
			return nil, errors.New("too big")
		}
	}
}

// AddBuffer appends data to buffer before sending it all-together.
func (c *ConnectionAdapter) AddBuffer(data ...interface{}) {
	c.bufout.Reset(c.conn)
	for _, buf := range data {
		switch value := buf.(type) {
		case string:
			logrus.Info("send: ", value)
			_, err := c.bufout.WriteString(value)
			if err != nil {
				logrus.Error(err)
			}
		case []byte:
			logrus.Info("send: ", string(value))
			_, err := c.bufout.Write(value)
			if err != nil {
				logrus.Error(err)
			}
		default:
			logrus.Errorf("failed to send session type %T", value)
			return
		}
		// Adds \r\n automatically at the end because smtp protocol requires it.
		_, _ = c.bufout.WriteString("\r\n")
	}
}

// Flush sends all buffered data at once for session client.
func (c *ConnectionAdapter) Flush() error {
	if c.bufout.Buffered() > 0 {
		c.resetTimeout()
		return c.bufout.Flush()
	}
	return nil
}

// SetDeadline set a timelimit for the connection to close automatically.
func (c *ConnectionAdapter) SetDeadline(d time.Duration) error {
	return c.conn.SetDeadline(time.Now().Add(d))
}

// RemoteAddr retrieves the session client real ip address.
func (c *ConnectionAdapter) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

// Close terminates the connection between session and server.
func (c *ConnectionAdapter) Close() error {
	return c.conn.Close()
}
