package smtp

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/jhillyerd/enmime"
	"github.com/sirupsen/logrus"
	"github.com/sonalys/letterme/domain/models"
)

const (
	greet = "220 %s Greetings"
	helo  = "250 %s Hello"

	ehlo               = "250 %s Hello\r\n"
	messageSize        = "250 SIZE %d\r\n"
	pipelining         = "250 PIPELINING\r\n"
	advertiseTLS       = "250 STARTTLS\r\n"
	enhancedStatusCode = "250 ENHANCEDSTATUSCODES\r\n"
	ok                 = "250 2.1.0 OK\r\n"
	rcvData            = "354 Ready\r\n"

	errNestedEmail  = "503 NESTEDEMAIL nested MAIL command"
	errTooBig       = "552 email is too big"
	errInvalidEmail = "450 4.7.1 %s"

	help = "250 HELP"
)

type command []byte

var (
	cmdHELO     command = []byte("HELO")
	cmdEHLO     command = []byte("EHLO")
	cmdHELP     command = []byte("HELP")
	cmdXCLIENT  command = []byte("XCLIENT")
	cmdMAIL     command = []byte("MAIL FROM:")
	cmdRCPT     command = []byte("RCPT TO:")
	cmdRSET     command = []byte("RSET")
	cmdVRFY     command = []byte("VRFY")
	cmdNOOP     command = []byte("NOOP")
	cmdQUIT     command = []byte("QUIT")
	cmdDATA     command = []byte("DATA")
	cmdSTARTTLS command = []byte("STARTTLS")
)

func (c command) match(in []byte) bool {
	lenIn := len(in)
	lenC := len(c)
	if lenIn < lenC {
		return false
	}
	return bytes.Equal(in[:lenC], c)
}

type clientState uint

const (
	clientStateGreeting clientState = iota
	clientStateCMD
	clientStateData
	clientStateShutdown
)

type Client struct {
	conn     *Connection
	envelope *models.UnencryptedEmail
	state    clientState
	tls      bool
	eSMTP    bool
}

func NewClient(conn *Connection) *Client {
	return &Client{
		conn: conn,
	}
}

func (c *Client) close() error {
	return c.conn.close()
}

func (c *Client) isAlive() bool {
	return c.conn != nil
}

const (
	_         = iota             // ignore first value by assigning to blank identifier
	KB uint32 = 1 << (10 * iota) // 1 << (10*1)
	MB                           // 1 << (10*2)
	GB                           // 1 << (10*3)
)

const maxDataSize = 25 * MB

func (c *Client) handle(s *Server) {
	for c.isAlive() {
		switch c.state {
		case clientStateGreeting:
			greeting := fmt.Sprintf(greet, s.c.Hostname)
			c.conn.send(greeting)
			c.state = clientStateCMD
		case clientStateCMD:
			c.parseCMD(s)
		case clientStateData:
			if uint32(c.conn.bufin.Buffered()) > maxDataSize {
				c.conn.send(errTooBig)
			}

			env, err := enmime.ReadEnvelope(c.conn.conn)
			if err != nil {
				c.conn.send(errInvalidEmail)
				c.resetTransaction()
				break
			}

			spew.Dump(env)
			c.resetTransaction()

		case clientStateShutdown:
			c.close()
		}

		if err := c.conn.flush(); err != nil {
			logrus.Error(err)
			return
		}
	}
}

func (c *Client) resetTransaction() {
	c.envelope = nil
}

func (c *Client) inTransaction() bool {
	return c.envelope != nil
}

func (c *Client) parseCMD(s *Server) {
	// TODO: fetch from above
	const maxLineSize = 4096
	line, err := c.conn.ReadLine()
	if err != nil {
		logrus.Error(line, err)
		c.close()
		return
	}

	logrus.Info("received: ", string(line))

	switch {
	case cmdHELO.match(line):
		c.resetTransaction()
		c.conn.send(fmt.Sprintf(helo, s.c.Hostname))
	case cmdEHLO.match(line):
		c.eSMTP = true
		c.resetTransaction()
		c.conn.send(
			fmt.Sprintf(ehlo, s.c.Hostname),
			fmt.Sprintf(messageSize, maxLineSize),
			pipelining,
			advertiseTLS,
			enhancedStatusCode,
			help,
		)
	case cmdHELP.match(line):
		c.conn.send(ok)
	case cmdMAIL.match(line):
		if c.inTransaction() {
			c.conn.send(errNestedEmail)
			break
		}

		buf := line[len(cmdMAIL):]
		addr, err := parseEmail(buf)
		if err != nil {
			c.conn.send(fmt.Sprintf(errInvalidEmail, err))
			break
		}

		c.envelope = &models.UnencryptedEmail{
			From: *addr,
		}
		c.conn.send(ok)
	case cmdRCPT.match(line):
		if len(c.envelope.ToList) > 10 {
			c.conn.send(fmt.Sprintf(errInvalidEmail, "too many recipients"))
			break
		}

		buf := line[len(cmdRCPT):]
		addr, err := parseEmail(buf)
		if err != nil {
			c.conn.send(fmt.Sprintf(errInvalidEmail, err))
			break
		}

		c.envelope.ToList = append(c.envelope.ToList, *addr)
		c.conn.send(ok)
	case cmdRSET.match(line):
		c.resetTransaction()
		c.conn.send(ok)
	case cmdQUIT.match(line):
		c.conn.send(ok)
		c.close()
	case cmdSTARTTLS.match(line):
		if err := c.conn.upgradeTLS(&s.pool.config.tlsConfig); err == nil {
			c.tls = true
		}
		c.resetTransaction()
	case cmdDATA.match(line):
		if len(c.envelope.ToList) == 0 {
			c.conn.send(errInvalidEmail)
			break
		}
		c.conn.send(rcvData)
		c.state = clientStateData
	case cmdNOOP.match(line):
		c.conn.send(ok)
	case cmdVRFY.match(line):
		c.conn.send(ok)
	}
}

func parseEmail(buf []byte) (*models.Address, error) {
	size := len(buf)
	if size < 4 || size > 253 {
		return nil, errors.New("address size must be between 4 and 253")
	}

	strip := buf[1 : len(buf)-1]
	addr := models.Address(strip)
	return &addr, addr.Validate()
}
