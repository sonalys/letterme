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
	AddSession(c net.Conn) (*Session, error)
	CloseSession(c *Session)
	Shutdown() <-chan bool
}

// SessionPool manages all the active sessions.
type SessionPool struct {
	wg             sync.WaitGroup
	sem            chan bool
	ctx            context.Context
	config         *PoolConfig
	activeSessions sync.Map
}

// PoolConfig are the configurations required to initialize the SessionPool
type PoolConfig struct {
	timeout   time.Duration
	capacity  uint
	hostname  string
	tlsConfig *tls.Config
}

// NewSessionPool instantiates a new connection pool.
func NewSessionPool(ctx context.Context, c *PoolConfig) *SessionPool {
	pool := &SessionPool{
		wg:             sync.WaitGroup{},
		sem:            make(chan bool, c.capacity),
		ctx:            ctx,
		config:         c,
		activeSessions: sync.Map{},
	}
	return pool
}

// addSession is when the semaphore allows one more session to the pool.
func (p *SessionPool) AddSession(c net.Conn) (*Session, error) {
	p.sem <- true
	p.wg.Add(1)
	sessionID := uuid.NewString()

	// As following SMTP, we don't want to start any connection already encrypted,
	// the client has to request the upgrade himself.
	conn, err := NewConnection(p.ctx, c, p.config.timeout, nil)
	if err != nil {
		return nil, err
	}

	client := NewSession(sessionID, conn)
	p.activeSessions.Store(sessionID, client)
	return client, nil
}

// closeSession is used when the connection is closed and the session is ending.
func (p *SessionPool) CloseSession(c *Session) {
	c.close()
	p.activeSessions.Delete(c.sessionID)
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
