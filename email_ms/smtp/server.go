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
	"github.com/sonalys/letterme/domain/cryptography"
	"github.com/sonalys/letterme/domain/utils"
)

// Server is the administrator for all receiving connections.
type Server struct {
	ctx      context.Context
	c        *ServerConfig
	pool     SessionManager
	listener net.Listener
}

// ServerConfigEnv defines the env name for server's configuration.
const ServerConfigEnv = "LM_SMTP_CONFIG"

// ServerConfig are the dependencies for the server to start.
type ServerConfig struct {
	MaxClients    uint                    `json:"max_clients"`
	MaxRecipients uint                    `json:"max_recipients"`
	MaxEmailSize  uint32                  `json:"max_email_size"`
	Timeout       time.Duration           `json:"timeout"`
	PrivateKey    cryptography.PrivateKey `json:"private_key"`
	Hostname      string                  `json:"hostname"`
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

	if c.PrivateKey.D == nil {
		errList = append(errList, newInvalidFieldErr("private_key"))
	}

	if c.Hostname == "" {
		errList = append(errList, newInvalidFieldErr("hostname"))
	}

	if len(errList) > 0 {
		return newInvalidConfigErr(errList)
	}
	return nil
}

// NewServer instantiates a new server.
func NewServer(ctx context.Context, c *ServerConfig) (*Server, error) {
	return &Server{
		c:   c,
		ctx: ctx,
		pool: NewSessionPool(ctx, &PoolConfig{
			tlsConfig: tls.Config{
				Rand:       rand.Reader,
				ClientAuth: tls.VerifyClientCertIfGiven,
				ServerName: c.Hostname,
				// TODO: fix this shit, should load from cert file.
				Certificates: []tls.Certificate{
					{PrivateKey: c.PrivateKey},
				},
			},
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

// Listen starts to accept new connections.
func (s *Server) Listen() error {
	addr := "localhost:2526"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to listen to %s", addr))
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
