package smtp

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

type Connection interface {
	UpgradeTLS(config *tls.Config) error
	ReadLine() ([]byte, error)
	ReadBytes(delim byte) ([]byte, error)
	ReadEnvelope(maxSize uint32) (reader io.Reader)
	SetDeadline(d time.Duration) error
	AddBuffer(data ...interface{})
	Flush() error
	RemoteAddr() string
	Close() error
}

type ConnectionAdapter struct {
	ctx      context.Context
	conn     net.Conn
	hostname string
	bufout   *bufio.Writer
	bufin    *bufio.Reader
	timeout  time.Duration
}

func NewConnection(ctx context.Context, c net.Conn, timeout time.Duration) Connection {
	return &ConnectionAdapter{
		ctx:     ctx,
		conn:    c,
		bufout:  bufio.NewWriter(c),
		bufin:   bufio.NewReader(c),
		timeout: timeout,
	}
}

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

func (c *ConnectionAdapter) ReadLine() ([]byte, error) {
	buf, err := c.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return bytes.Trim(buf, "\r\n"), nil
}

func (c *ConnectionAdapter) ReadBytes(delim byte) ([]byte, error) {
	buf, err := c.bufin.ReadBytes(delim)
	if err != nil {
		return nil, err
	}
	c.resetTimeout()
	return buf, nil
}

func (c *ConnectionAdapter) ReadEnvelope(maxSize uint32) (reader io.Reader) {
	var buffer bytes.Buffer
	var size uint32
	go func() {
		for {
			slice, err := c.ReadLine()
			if err != nil {
				c.Close()
				return
			}

			size = size + uint32(len(slice))
			if size > maxSize {
				break
			}

			if bytes.Equal(slice, []byte{'.', '\r'}) {
				break
			}
			buffer.Write(slice)
		}
	}()
	return &buffer
}

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
			_, err := c.bufout.Write(value)
			if err != nil {
				logrus.Error(err)
			}
		}
	}

	_, err := c.bufout.WriteString("\r\n")
	if err != nil {
		logrus.Error(err)
	}
}

func (c *ConnectionAdapter) Flush() error {
	if c.bufout.Buffered() > 0 {
		c.resetTimeout()
		return c.bufout.Flush()
	}
	return nil
}

func (c *ConnectionAdapter) SetDeadline(d time.Duration) error {
	return c.conn.SetDeadline(time.Now().Add(d))
}

func (c *ConnectionAdapter) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *ConnectionAdapter) Close() error {
	return c.conn.Close()
}
