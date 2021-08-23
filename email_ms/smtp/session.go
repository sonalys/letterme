package smtp

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/jhillyerd/enmime"
	"github.com/sirupsen/logrus"
	"github.com/sonalys/letterme/domain/models"
)

const (
	// Handshake section.

	// HELO protocol uses CODE\sMESSAGE format.
	messageGreet = "220 %s Greetings"
	messageHELO  = "250 %s Hello"

	// // // // // // // // //

	// Commands section.

	// EHLO protocol needs CODE-COMMAND\r\n format.
	messageEHLO         = "250-%s Hello\r\n"
	messageAdvertiseTLS = "250-STARTTLS\r\n"
	messageSize         = "250-SIZE %d\r\n"
	messagePipelining   = "250-PIPELINING\r\n"
	messageEnhance      = "250-ENHANCEDSTATUSCODES\r\n"
	// Last command doesn't need \r\n because it's inserted automatically.
	messageHelp = "250 HELP"

	// // // // // // // // //

	// Response section.

	responseOK    = "250 2.1.0 OK"
	responseReady = "354 Ready"

	// // // // // // // // //

	// Error section.

	errNestedEmail  = "503 NESTEDEMAIL nested MAIL command\r\n"
	errTooBig       = "552 email is too big\r\n"
	errInvalidEmail = "450 4.7.1 %s\r\n"
)

const (
	_         = iota             // ignore first value by assigning to blank identifier
	KB uint32 = 1 << (10 * iota) // 1 << (10*1)
	MB                           // 1 << (10*2)
	GB                           // 1 << (10*3)
)

type command []byte

var (
	cmdHELO     command = []byte("HELO")
	cmdEHLO     command = []byte("EHLO")
	cmdHELP     command = []byte("HELP")
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

// Session is created when the server has accepted the connection.
// TODO: needs to create more security logic here, reputation, session timeout...
type Session struct {
	// sessionID is a unique identifier for this session.
	sessionID string
	// conn is the connection handler betwen session and client.
	conn Connection
	// envelope is the current email transaction, if nil, the session is not creating an email.
	envelope *models.UnencryptedEmail
	// state holds data for the smtp state machine.
	state clientState
	// tls is true when the session upgraded the conn to tls.
	tls bool
	// eSMTP is true when client used EHLO protocol.
	eSMTP bool
}

// NewSession instantiates a new session
func NewSession(sessionID string, conn Connection) *Session {
	session := &Session{
		// TODO: maybe fetch this from account_ms.
		sessionID: sessionID,
		conn:      conn,
	}
	return session
}

// close stops the session.
func (c *Session) close() error {
	err := c.conn.Close()
	c.conn = nil
	return err
}

// isAlive check if the connection is still active,
// critical for leaving session routine.
func (c *Session) isAlive() bool {
	return c.conn != nil
}

// startSession creates a new routine for handling client inputs.
// it stops when the session timeouts, gives invalid input or closes.
func (c *Session) startSession(s *Server) {
	for c.isAlive() {
		switch c.state {
		case clientStateGreeting:
			c.state = clientStateCMD
			c.conn.AddBuffer(fmt.Sprintf(messageGreet, s.c.Hostname))
		case clientStateCMD:
			c.parseCMD(s)
		case clientStateData:
			if err := c.readData(); err != nil {
				c.conn.AddBuffer(fmt.Sprintf(errInvalidEmail, err))
			} else {
				c.conn.AddBuffer(responseOK)
			}
			c.resetTransaction()

		case clientStateShutdown:
			return
		}

		if err := c.conn.Flush(); err != nil {
			logrus.Error(err)
			return
		}
	}
}

// readData is a handler for the clientData state.
// it parses the envelope from the client.
func (c *Session) readData() error {
	const maxSize = 25 * MB
	_, err := enmime.ReadEnvelope(c.conn.ReadEnvelope(maxSize))
	if err != nil {
		c.conn.AddBuffer(fmt.Sprintf(errInvalidEmail, err))
		c.resetTransaction()
		return err
	}
	// Do something with the envelope ;)
	return nil
}

// resetTransaction discards all buffered data from email,
// it remains in CMD state.
func (c *Session) resetTransaction() {
	c.envelope = nil
}

// inTransaction indicates if the current session is creating an email.
func (c *Session) inTransaction() bool {
	return c.envelope != nil
}

// parseCMD is a parser for the clientCMD state, it reads the client inputs and process it properly.
func (c *Session) parseCMD(s *Server) {
	// value extracted from RFC-5321.
	const maxLineSize = 512
	line, err := c.conn.ReadLine()
	if err != nil {
		logrus.Error(line, err)
		c.close()
		return
	}

	switch {
	case cmdHELO.match(line):
		c.resetTransaction()
		c.conn.AddBuffer(fmt.Sprintf(messageHELO, s.c.Hostname))
	case cmdEHLO.match(line):
		c.eSMTP = true
		c.resetTransaction()
		c.conn.AddBuffer(
			fmt.Sprintf(messageEHLO, s.c.Hostname),
			fmt.Sprintf(messageSize, maxLineSize),
			messagePipelining,
			// advertiseTLS, // disabled in debug because we don't have any certificate
			messageEnhance,
			messageHelp,
		)
	case cmdHELP.match(line):
		c.conn.AddBuffer(responseOK)
	case cmdMAIL.match(line):
		if c.inTransaction() {
			c.conn.AddBuffer(errNestedEmail)
			break
		}

		buf := line[len(cmdMAIL):]
		addr, err := parseEmail(buf)
		if err != nil {
			c.conn.AddBuffer(fmt.Sprintf(errInvalidEmail, err))
			break
		}

		c.envelope = &models.UnencryptedEmail{
			From: *addr,
		}
		c.conn.AddBuffer(responseOK)
	case cmdRCPT.match(line):
		const maxRecipients = 100
		if len(c.envelope.ToList) > maxRecipients {
			c.conn.AddBuffer(fmt.Sprintf(errInvalidEmail, "too many recipients"))
			break
		}

		buf := line[len(cmdRCPT):]
		addr, err := parseEmail(buf)
		if err != nil {
			c.conn.AddBuffer(fmt.Sprintf(errInvalidEmail, err))
			break
		}

		c.envelope.ToList = append(c.envelope.ToList, *addr)
		c.conn.AddBuffer(responseOK)
	case cmdRSET.match(line):
		c.resetTransaction()
		c.conn.AddBuffer(responseOK)
	case cmdQUIT.match(line):
		c.conn.AddBuffer(responseOK)
		c.close()
	case cmdSTARTTLS.match(line):
		// TODO: fix this, should fetch cert from server config.
		// We need to check if the client also uses a digital certificate, we don't want to receive emails from untrusted entities.
		if err := c.conn.UpgradeTLS(nil); err == nil {
			c.tls = true
		} else {
			logrus.Info(err)
		}
		c.resetTransaction()
	case cmdDATA.match(line):
		if len(c.envelope.ToList) == 0 {
			c.conn.AddBuffer(fmt.Sprintf(errInvalidEmail, "no recipients"))
			break
		}
		c.conn.AddBuffer(responseReady)
		c.state = clientStateData
	case cmdNOOP.match(line):
		c.conn.AddBuffer(responseOK)
	case cmdVRFY.match(line):
		// We don't reveal what addresses we have or not, for privacy reasons.
		c.conn.AddBuffer(responseOK)
	}
}

// parseEmail parses <email@example.com> to models.Address.
//
// TODO: maybe should be in domain?
func parseEmail(buf []byte) (*models.Address, error) {
	size := len(buf)
	if size < 4 || size > 253 {
		return nil, errors.New("address size must be between 4 and 253")
	}

	strip := buf[1 : len(buf)-1]
	addr := models.Address(strip)
	return &addr, addr.Validate()
}
