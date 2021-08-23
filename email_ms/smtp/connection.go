package smtp

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

type Connection struct {
	ctx     context.Context
	conn    net.Conn
	bufout  *bufio.Writer
	bufin   *bufio.Reader
	timeout time.Duration
}

func NewConnection(ctx context.Context, c net.Conn, timeout time.Duration) *Connection {
	return &Connection{
		ctx:     ctx,
		conn:    c,
		bufout:  bufio.NewWriter(c),
		bufin:   bufio.NewReader(c),
		timeout: timeout,
	}
}

func (c *Connection) upgradeTLS(config *tls.Config) error {
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

func (c *Connection) resetTimeout() {
	// we don't wan't to keep refreshing this deadline when the ms is closing.
	if c.ctx.Err() == nil {
		c.setDeadline(c.timeout)
	}
}

func (c *Connection) ReadLine() ([]byte, error) {
	buf, err := c.bufin.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return bytes.Trim(buf, "\r\n"), nil
}

func (c *Connection) send(data ...interface{}) {
	c.bufout.Reset(c.conn)
	for _, buf := range data {
		switch value := buf.(type) {
		case string:
			logrus.Info("Sending: ", value)
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

func (c *Connection) flush() error {
	if c.bufout.Buffered() > 0 {
		return c.bufout.Flush()
	}
	return nil
}

func (c *Connection) setDeadline(d time.Duration) error {
	return c.conn.SetDeadline(time.Now().Add(d))
}

func (c *Connection) remoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Connection) close() error {
	return c.conn.Close()
}
