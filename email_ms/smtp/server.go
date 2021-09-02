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
	"github.com/sonalys/letterme/domain/utils"
)

// Server is the administrator for all receiving connections.
type Server struct {
	ctx      context.Context
	c        *ServerConfig
	tls      *tls.Config
	pool     SessionManager
	listener net.Listener
}

// ServerConfigEnv defines the env name for server's configuration.
const ServerConfigEnv = "LM_SMTP_CONFIG"

// ServerConfig are the dependencies for the server to start.
type ServerConfig struct {
	MaxClients     uint          `json:"max_clients"`
	MaxRecipients  uint          `json:"max_recipients"`
	MaxEmailSize   uint32        `json:"max_email_size"`
	Timeout        time.Duration `json:"timeout"`
	CertificateKey string        `json:"certificate_key"`
	Certificate    string        `json:"certificate"`
	Hostname       string        `json:"hostname"`
	Address        string        `json:"address"`
}

// Validate implements the validatable interface.
func (c ServerConfig) Validate() error {
	var errList []error

	if c.MaxClients == 0 {
		errList = append(errList, newInvalidFieldErr("max_clients"))
	}

	if c.Timeout == 0 {
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

	return &Server{
		c:   c,
		ctx: ctx,
		tls: tlsConfig,
		pool: NewSessionPool(ctx, &PoolConfig{
			capacity: c.MaxClients,
			timeout:  c.Timeout,
		}),
	}, nil
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
	// TODO: improve this
	certificate, err := tls.X509KeyPair([]byte(c.Certificate), []byte(c.CertificateKey))
	if err != nil {
		return nil, newInvalidCertificateErr(err)
	}

	return &tls.Config{
		Rand:         rand.Reader,
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.VerifyClientCertIfGiven,
		ServerName:   c.Hostname,
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
		s.pool.HandleConnection(conn, s)
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
