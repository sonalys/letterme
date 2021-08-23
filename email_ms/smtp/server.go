// nolint // fuck you
package smtp

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"net"
	"time"

	"github.com/sonalys/letterme/domain/cryptography"
	"github.com/sonalys/letterme/domain/utils"
)

type serverState uint

const (
	serverStateNew serverState = iota
	serverStateStopped
	serverStateRunning
	serverStateError
)

type Server struct {
	ctx      context.Context
	c        *ServerConfig
	state    serverState
	pool     *ConnectionPool
	listener net.Listener
}

const serverConfigEnv = "LM_SMTP_CONFIG"

type ServerConfig struct {
	MaxClients uint                    `json:"max_clients"`
	Timeout    time.Duration           `json:"timeout"`
	PrivateKey cryptography.PrivateKey `json:"private_key"`
	Hostname   string                  `json:"hostname"`
}

func (c ServerConfig) Validate() error {
	if c.MaxClients == 0 {
		return errors.New("invalid max_clients value")
	}

	if c.Timeout == 0 {
		return errors.New("invalid timeout value")
	}

	if c.PrivateKey.D == nil {
		return errors.New("invalid private_key value")
	}

	if c.Hostname == "" {
		return errors.New("invalid hostname")
	}
	return nil
}

func NewServer(ctx context.Context, c *ServerConfig) (*Server, error) {
	return &Server{
		c:     c,
		ctx:   ctx,
		state: serverStateNew,
		// TODO: use env to config this
		pool: NewConnectionPool(ctx, &PoolConfig{
			tlsConfig: tls.Config{
				Rand:       rand.Reader,
				ClientAuth: tls.VerifyClientCertIfGiven,
				ServerName: c.Hostname,
				Certificates: []tls.Certificate{
					{PrivateKey: c.PrivateKey},
				},
			},
			capacity: c.MaxClients,
			timeout:  c.Timeout,
		}),
	}, nil
}

func InitServerFromEnv(ctx context.Context) (*Server, error) {
	cfg := new(ServerConfig)
	if err := utils.LoadFromEnv(serverConfigEnv, cfg); err != nil {
		return nil, err
	}
	return NewServer(ctx, cfg)
}

func (s *Server) Listen() error {
	listener, err := net.Listen("tcp", "localhost:2526")
	if err != nil {
		return err
	}

	s.listener = listener
	s.state = serverStateRunning

	for {
		conn, err := listener.Accept()
		if err != nil {
			if err, ok := err.(net.Error); ok && !err.Temporary() {
				s.pool.Shutdown()
				s.state = serverStateStopped
				return nil
			}
			// log error
			continue
		}
		s.pool.handleConnection(conn, s)
	}
}

func (s *Server) Shutdown() {
	if s.listener != nil {
		_ = s.listener.Close()
		<-s.pool.Shutdown()
	}
	s.state = serverStateStopped
}
