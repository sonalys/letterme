package smtp

import (
	"context"
	"crypto/tls"
	"net"
	"sync"
	"time"
)

type ConnectionPool struct {
	ctx           context.Context
	sem           chan bool
	activeClients sync.Map
	wg            sync.WaitGroup
	config        *PoolConfig
}

type PoolConfig struct {
	tlsConfig tls.Config
	capacity  uint
	timeout   time.Duration
}

func NewConnectionPool(ctx context.Context, c *PoolConfig) *ConnectionPool {
	pool := &ConnectionPool{
		ctx:           ctx,
		config:        c,
		sem:           make(chan bool, c.capacity),
		activeClients: sync.Map{},
		wg:            sync.WaitGroup{},
	}
	return pool
}

func (p *ConnectionPool) handleConnection(c net.Conn, s *Server) {
	client, err := p.addConnection(c)
	if err != nil {
		return
	}
	go func() {
		client.handle(s)
		p.removeClient(client)
	}()
}

func (p *ConnectionPool) addConnection(c net.Conn) (*Client, error) {
	p.sem <- true
	p.wg.Add(1)
	client := NewClient(NewConnection(p.ctx, c, p.config.timeout))
	p.activeClients.Store(client.conn.remoteAddr(), client)
	return client, nil
}

func (p *ConnectionPool) removeClient(c *Client) {
	c.close()
	p.activeClients.Delete(c.conn.remoteAddr())
	p.wg.Done()
	<-p.sem
}

func (p *ConnectionPool) Shutdown() <-chan bool {
	timeOut := time.Second
	p.activeClients.Range(func(key, value interface{}) bool {
		client, ok := value.(*Client)
		if !ok {
			return true
		}

		_ = client.conn.setDeadline(timeOut)
		return true
	})

	closeChan := make(chan bool, 1)
	go func() {
		p.wg.Wait()
		closeChan <- true
	}()
	return closeChan
}
