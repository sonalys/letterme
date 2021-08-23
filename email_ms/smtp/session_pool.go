package smtp

import (
	"context"
	"crypto/tls"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
)

// SessionManager is the signature required to manage all server sessions.
type SessionManager interface {
	HandleConnection(c net.Conn, s *Server)
	Shutdown() <-chan bool
}

// SessionPool manages all the active sessions.
type SessionPool struct {
	ctx context.Context
	// sem is a semaphore used to control max active sessions.
	sem chan bool
	// activeSessions are all the active sessions.
	activeSessions sync.Map
	// wg is used to wait for all the sessions to finish before shutting down.
	wg sync.WaitGroup
	// config is used to hold all config data.
	config *PoolConfig
}

// PoolConfig are the configurations required to initialize the SessionPool
type PoolConfig struct {
	tlsConfig tls.Config
	capacity  uint
	timeout   time.Duration
}

// NewSessionPool instantiates a new connection pool.
func NewSessionPool(ctx context.Context, c *PoolConfig) *SessionPool {
	pool := &SessionPool{
		ctx:            ctx,
		config:         c,
		sem:            make(chan bool, c.capacity),
		activeSessions: sync.Map{},
		wg:             sync.WaitGroup{},
	}
	return pool
}

// HandleConnection receives a TCP connection and manages a new session.
func (p *SessionPool) HandleConnection(c net.Conn, s *Server) {
	p.sem <- true
	sessionID := uuid.NewString()
	client := p.addSession(c, sessionID)
	go func() {
		client.startSession(s)
		p.closeSession(client, sessionID)
	}()
}

// addSession is when the semaphore allows one more session to the pool.
func (p *SessionPool) addSession(c net.Conn, key string) *Session {
	p.wg.Add(1)
	client := NewSession(key, NewConnection(p.ctx, c, p.config.timeout))
	p.activeSessions.Store(key, client)
	return client
}

// closeSession is used when the connection is closed and the session is ending.
func (p *SessionPool) closeSession(c *Session, key string) {
	c.close()
	p.activeSessions.Delete(key)
	p.wg.Done()
	<-p.sem
}

// Shutdown is used when the listener is closed, and we need to wait for lasting sessions.
// We give them 10 seconds before brute-force closing the TCP connection.
func (p *SessionPool) Shutdown() <-chan bool {
	// TODO: use config to set this value.
	timeOut := 10 * time.Second
	p.activeSessions.Range(func(key, value interface{}) bool {
		client, ok := value.(*Session)
		if !ok {
			return true
		}

		_ = client.conn.SetDeadline(timeOut)
		return true
	})

	// create a close channel for the server to wait for this function to finish.
	closeChan := make(chan bool, 1)
	go func() {
		p.wg.Wait()
		closeChan <- true
	}()
	return closeChan
}
