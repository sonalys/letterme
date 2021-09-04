// nolint // fuck you
package smtp

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sonalys/letterme/domain/models"
	"github.com/sonalys/letterme/domain/utils"
)

// Server is the administrator for all receiving connections.
type Server struct {
	c              *ServerConfig
	tls            *tls.Config
	ctx            context.Context
	pool           SessionManager
	listener       net.Listener
	cachedMessages [][]byte
	pipeline       Pipeline
}

// ServerConfigEnv defines the env name for server's configuration.
const ServerConfigEnv = "LM_SMTP_CONFIG"

// ServerConfig are the dependencies for the server to start.
type ServerConfig struct {
	Address        string        `json:"address"`
	Hostname       string        `json:"hostname"`
	MaxClients     uint          `json:"max_clients"`
	Certificate    string        `json:"certificate"`
	MaxEmailSize   uint32        `json:"max_email_size"`
	MaxRecipients  uint          `json:"max_recipients"`
	CertificateKey string        `json:"certificate_key"`
	SessionTimeout time.Duration `json:"timeout"`
}

// Validate implements the validatable interface.
func (c ServerConfig) Validate() error {
	var errList []error

	if c.MaxClients == 0 {
		errList = append(errList, newInvalidFieldErr("max_clients"))
	}

	if c.SessionTimeout == 0 {
		errList = append(errList, newInvalidFieldErr("timeout"))
	}

	if len(c.CertificateKey) == 0 {
		errList = append(errList, newInvalidFieldErr("certificate_key"))
	}

	if len(c.Certificate) == 0 {
		errList = append(errList, newInvalidFieldErr("certificate"))
	}

	if c.Hostname == "" {
		errList = append(errList, newInvalidFieldErr("hostname"))
	}

	if c.Address == "" {
		errList = append(errList, newInvalidFieldErr("address"))
	}

	if len(errList) > 0 {
		return newInvalidConfigErr(errList)
	}
	return nil
}

// NewServer instantiates a new server.
func NewServer(ctx context.Context, c *ServerConfig) (*Server, error) {
	tlsConfig, err := loadTLS(c)
	if err != nil {
		return nil, newInitializeServerErr(err)
	}

	server := &Server{
		c:        c,
		ctx:      ctx,
		tls:      tlsConfig,
		pipeline: Pipeline{},
		pool: NewSessionPool(ctx, &PoolConfig{
			hostname:  c.Hostname,
			capacity:  c.MaxClients,
			timeout:   c.SessionTimeout,
			tlsConfig: tlsConfig,
		}),
	}

	server.cacheMessages()
	return server, nil
}

// InitServerFromEnv loads all the configs from env and starts the server.
func InitServerFromEnv(ctx context.Context) (*Server, error) {
	cfg := new(ServerConfig)
	if err := utils.LoadFromEnv(ServerConfigEnv, cfg); err != nil {
		return nil, errors.Wrap(err, "failed to initialize server")
	}
	return NewServer(ctx, cfg)
}

func loadTLS(c *ServerConfig) (*tls.Config, error) {
	certificate, err := tls.X509KeyPair([]byte(c.Certificate), []byte(c.CertificateKey))
	if err != nil {
		return nil, newInvalidCertificateErr(err)
	}

	return &tls.Config{
		Rand:         rand.Reader,
		Certificates: []tls.Certificate{certificate},
		// If you change this I will fucking kill you,
		// we need to be sure that every provider using this SMTP is verified.
		ClientAuth: tls.RequireAndVerifyClientCert,
	}, nil
}

// Listen starts to accept new connections.
func (s *Server) Listen() error {
	listener, err := net.Listen("tcp", s.c.Address)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to listen to %s", s.c.Address))
	}

	s.listener = listener

	for {
		// this will return err imediately after close method is called.
		conn, err := listener.Accept()
		if err != nil {
			// this indicates that the listener was closed.
			if err, ok := err.(net.Error); ok && !err.Temporary() {
				s.pool.Shutdown()
				return nil
			}
			logrus.Error("failed to accept new tcp connection", err)
			continue
		}
		// Session pool ensures there's a slot for this connection.
		session, err := s.pool.AddSession(conn)
		if err != nil {
			logrus.Error("failed to create session: ", err)
		} else {
			// After ensuring the slot, we can create a routine to handle it.
			go s.handleSession(session)
		}
	}
}

// Shutdown tells the server that is should no longer accept new connections,
// and wait for all existent sessions to finish before closing.
func (s *Server) Shutdown() {
	if s.listener != nil {
		_ = s.listener.Close()
		<-s.pool.Shutdown()
	}
}

// AddMiddleware adds a new middleware to the envelope pipeline.
func (s *Server) AddMiddlewares(middleware ...EnvelopeMiddleware) {
	middlewares := append(middleware, releaseEnvelopeMiddleware)
	s.pipeline.AddMiddlewares(middlewares...)
}

// releaseEnvelopeMiddleware releases the envelope from the heap pool.
func releaseEnvelopeMiddleware(next EnvelopeHandler) EnvelopeHandler {
	return func(envelope *models.UnencryptedEmail) error {
		envelopePool.Put(envelope)
		return nil
	}
}

// startEnvelopePipeline process all the handlers for this envelope
func (s *Server) startEnvelopePipeline(envelope *models.UnencryptedEmail) {
	if err := s.pipeline.Start(envelope); err != nil {
		logrus.Error(err)
	}
}

func (s *Server) handleSession(session *Session) {
	defer s.pool.CloseSession(session)
	cache := s.cachedMessages

	for session.isAlive() {
		switch session.state {
		case clientStateGreeting:
			session.state = clientStateCMD
			session.Send(cache[mGreet])
		case clientStateCMD:
			s.handleCommand(session)
		case clientStateData:
			if err := session.readEnvelope(); err != nil {
				logrus.Error("failed to parse session data", err)
				session.Send(cache[mErrInvalidEmail])
				return
			}

			go s.startEnvelopePipeline(session.envelope)

			session.state = clientStateCMD
			session.resetTransaction()
			session.Send(cache[mOK])
		}
	}
}

func (s *Server) handleCommand(session *Session) {
	cache := s.cachedMessages
	line, err := session.ReadLine()
	if err != nil {
		session.close()
		return
	}

	switch {
	case cmdHELO.match(line):
		session.resetTransaction()
		session.Send(cache[mHELO])
	case cmdEHLO.match(line):
		session.resetTransaction()
		session.Send(
			cache[mEHLO],
			cache[mSize],
			cache[mPipeline],
			cache[mTLS],
			cache[mEnhance],
			cache[mHelp],
		)
	case cmdHELP.match(line):
		session.Send(cache[mOK])
	case cmdMAIL.match(line):
		if session.inTransaction() {
			session.Send(cache[mErrTransaction])
			break
		}
		buf := line[len(cmdMAIL):]
		addr, err := parseEmailAddress(buf)
		if err != nil {
			session.Send(cache[mErrInvalidEmail])
			break
		}

		envelope := envelopePool.Get().(*models.UnencryptedEmail)
		envelope.From = *addr

		session.envelope = envelope
		session.Send(cache[mOK])
	case cmdRCPT.match(line):
		const maxRecipients = 100
		if len(session.envelope.ToList) > maxRecipients {
			session.Send(cache[mErrTooManyRecipients])
			break
		}

		buf := line[len(cmdRCPT):]
		addr, err := parseEmailAddress(buf)
		if err != nil {
			session.Send(cache[mErrInvalidEmail])
			break
		}

		if addr.Domain() != s.c.Hostname {
			session.Send(cache[mErrInvalidEmail])
			break
		}

		session.envelope.ToList = append(session.envelope.ToList, *addr)
		session.Send(cache[mOK])
	case cmdRSET.match(line):
		session.resetTransaction()
		session.Send(cache[mOK])
	case cmdQUIT.match(line):
		session.Send(cache[mOK])
		session.close()
	case cmdSTARTTLS.match(line):
		// TODO: fix this, should fetch cert from server config.
		// We need to check if the client also uses a digital certificate, we don't want to receive emails from untrusted entities.
		if err := session.conn.UpgradeTLS(s.tls); err == nil {
			session.tls = true
		} else {
			logrus.Info(err)
		}
		session.resetTransaction()
	case cmdDATA.match(line):
		if !session.inTransaction() {
			session.Send(cache[mErrTransaction])
			break
		}
		if len(session.envelope.ToList) == 0 {
			session.Send(cache[mErrNoRecipients])
			break
		}
		session.Send(cache[mReady])
		session.state = clientStateData
	case cmdNOOP.match(line):
		session.Send(cache[mOK])
	case cmdVRFY.match(line):
		// We don't reveal what addresses we have or not, for privacy reasons.
		session.Send(cache[mOK])
	}
}

// cacheMessages cache all messages so we don't have to build strings during runtime.
func (s *Server) cacheMessages() {
	s.cachedMessages = [][]byte{
		[]byte(fmt.Sprintf("220 %s Greetings", s.c.Hostname)),
		[]byte(fmt.Sprintf("250 %s Hello", s.c.Hostname)),
		[]byte(fmt.Sprintf("250-%s Hello", s.c.Hostname)),
		[]byte("250-STARTTLS"),
		[]byte(fmt.Sprintf("250-SIZE %d", maxEnvelopeDataSize)),
		[]byte("250-PIPELINING"),
		[]byte("250-ENHANCEDSTATUSCODES"),
		[]byte("250 HELP"),
		[]byte("250 OK"),
		[]byte("354 READY"),
		[]byte("503 ALREADY IN TRANSACTION"),
		[]byte("552 EMAIL IS TOO BIG"),
		[]byte("552 EMAIL IS INVALID"),
		[]byte("552 NO RECIPIENTS"),
		[]byte("552 TOO MANY RECIPIENTS"),
	}
}
