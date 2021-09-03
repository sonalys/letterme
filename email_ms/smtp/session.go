package smtp

import (
	"sync"

	"github.com/jhillyerd/enmime"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sonalys/letterme/domain/models"
)

type clientState uint

const (
	clientStateGreeting clientState = iota
	clientStateCMD
	clientStateData
)

// envelopePool is a buffer pool to avoid memory allocation.
var envelopePool = sync.Pool{
	New: func() interface{} {
		return models.NewUnencryptedEmail()
	},
}

// Session is created when the server has accepted the connection.
//
// Session is responsible for holding the current state for the smtp protocol,
// for each individual connection.
//
// TODO: needs to create more security logic here, reputation, session timeout...
type Session struct {
	sessionID string
	conn      Connection
	envelope  *models.UnencryptedEmail
	state     clientState
	tls       bool
	closed    bool
}

// NewSession instantiates a new session.
func NewSession(sessionID string, conn Connection) *Session {
	session := &Session{
		// TODO: maybe fetch this from account_ms.
		sessionID: sessionID,
		conn:      conn,
	}
	return session
}

// Send is used to send data to client.
func (c *Session) Send(data ...interface{}) error {
	c.conn.AddBuffer(data...)
	return c.conn.Flush()
}

func (c *Session) ReadLine() ([]byte, error) {
	line, err := c.conn.ReadLine()
	logrus.Info("received: ", string(line))
	return line, err
}

// close stops the session.
func (c *Session) close() error {
	c.closed = true
	err := c.conn.Close()
	return errors.Wrap(err, "failed to close session connection")
}

// isAlive check if the connection is still active,
// critical for leaving session routine.
func (c *Session) isAlive() bool {
	return !c.closed
}

// readEnvelope is a handler for the clientData state.
// it parses the envelope from the client.
func (c *Session) readEnvelope() error {
	buf, err := c.conn.ReadEnvelope()
	if err != nil {
		logrus.Infof("error receiving envelope: %s", err)
		return err
	}

	if buf == nil {
		return nil
	}

	env, err := enmime.ReadEnvelope(buf)
	if err != nil {
		return errors.Wrap(err, "failed to parse data")
	}

	if len(env.Errors) > 0 {
		return errors.New("invalid envelope")
	}

	c.envelope.Text = []byte(env.Text)
	c.envelope.HTML = []byte(env.HTML)

	for _, attachment := range env.Attachments {
		c.envelope.Attachments = append(c.envelope.Attachments, models.AttachmentRequest{
			Buffer:   attachment.Content,
			Filename: attachment.FileName,
			Attachment: models.Attachment{
				Size:        uint32(len(attachment.Content)),
				MimeType:    models.MimeType(attachment.ContentType),
				ContentID:   attachment.ContentID,
				Disposition: attachment.Disposition,
				Insecure:    true,
			},
		})
	}

	for _, inline := range env.Inlines {
		c.envelope.Inlines = append(c.envelope.Inlines, models.AttachmentRequest{
			Buffer:   inline.Content,
			Filename: inline.FileName,
			Attachment: models.Attachment{
				Size:        uint32(len(inline.Content)),
				MimeType:    models.MimeType(inline.ContentType),
				ContentID:   inline.ContentID,
				Disposition: inline.Disposition,
				Insecure:    true,
			},
		})
	}

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
