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
}

// NewConnection instantiates a new connection controller.
func NewConnection(ctx context.Context, c net.Conn, timeout time.Duration, tls *tls.Config) Connection {
	conn := &ConnectionAdapter{
		ctx:     ctx,
		conn:    c,
		bufout:  bufio.NewWriter(c),
		bufin:   bufio.NewReader(c),
		timeout: timeout,
	}
	// try to upgrade to tls, doesn't matter if it fails.
	conn.UpgradeTLS(tls)
	return conn
}

// UpgradeTLS upgrades the socket connection with tls encryption.
func (c *ConnectionAdapter) UpgradeTLS(config *tls.Config) error {
	// wrap c.conn in a new TLS server side connection
	tlsConn := tls.Server(c.conn, config)

	// Call handshake here to get any handshake error before reading starts
	err := tlsConn.Handshake()
	if err != nil {
		return err
	}

	c.conn = tlsConn
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
	return bytes.Trim(buf, "\r\n"), nil
}

// ReadBytes reads the buffer until it reaches the first 'delim' byte.
func (c *ConnectionAdapter) ReadBytes(delim byte) ([]byte, error) {
	buf, err := c.bufin.ReadBytes(delim)
	if err != nil {
		return nil, err
	}
	c.resetTimeout()
	return buf, nil
}

// maxEnvelopeDataSize oh yes, hardcoded value, 25 MB per envelope buffer
const maxEnvelopeDataSize = 25 * MB

// bufferPool pre-allocates buffers for parsing envelopes.
var bufferPool = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, KB)
		return &buf
	},
}

// outputPool pre-allocates buffers for parsing envelopes.
var outputPool = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 0, 25*MB)
		return &buf
	},
}

// ReadEnvelope parses the upcoming data buffer.
func (c *ConnectionAdapter) ReadEnvelope() (reader io.Reader, err error) {
	// We allocate 1 buffer of 25 MB for the whole envelope
	// We allocate 1 buffer of 1 KB to gradually read the tcp conn.
	// TODO: use only 1 buffer to do it all.
	output := bytes.NewBuffer(*outputPool.Get().(*[]byte))
	buf := bufferPool.Get().(*[]byte)
	var size uint32
	defer bufferPool.Put(buf)

	for {
		n, err := c.conn.Read(*buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return reader, nil
			}
			c.Close()
			return reader, err
		}

		size = size + uint32(n)
		if size > maxEnvelopeDataSize {
			return reader, errors.New("too big")
		}

		if n > 3 && bytes.Equal((*buf)[n-3:n-1], []byte{'.', '\r'}) {
			return reader, nil
		}
		output.Write((*buf)[:n])
	}
}

// AddBuffer appends data to buffer before sending it all-together.
func (c *ConnectionAdapter) AddBuffer(data ...interface{}) {
	c.bufout.Reset(c.conn)
	for _, buf := range data {
		switch value := buf.(type) {
		case string:
			_, err := c.bufout.WriteString(value)
			if err != nil {
				logrus.Error(err)
			}
		case []byte:
			_, err := c.bufout.Write(value)
			if err != nil {
				logrus.Error(err)
			}
		default:
			logrus.Errorf("failed to send session type %T", value)
		}
	}

	_, _ = c.bufout.WriteString("\r\n")
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
