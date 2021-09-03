package smtp

import (
	"context"
	"crypto/tls"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/sonalys/letterme/domain/utils"
	"github.com/stretchr/testify/require"
)

func Test_NewConnection(t *testing.T) {
	ctx := context.Background()
	conn := NewConnection(ctx, &net.TCPConn{}, time.Hour, nil)
	require.NotNil(t, conn)
}

func Test_ConnectionUpgradeTLS(t *testing.T) {
	cfg := new(ServerConfig)
	err := utils.LoadFromEnv(ServerConfigEnv, cfg)
	require.NoError(t, err, "should load config from env")

	tlsConfig, err := loadTLS(cfg)
	require.NoError(t, err, "should initialize server config")
	// We won't use real certificate here.
	tlsConfig.InsecureSkipVerify = true
	// We can't verify self-signed client certificates here.
	tlsConfig.ClientAuth = tls.RequireAnyClientCert

	sv, err := net.Listen("tcp", ":1234")
	require.NoError(t, err, "should initialize tls server")
	defer sv.Close()

	var wg sync.WaitGroup

	go func() {
		defer wg.Done()
		wg.Add(1)
		conn, err := sv.Accept()
		require.NoError(t, err)

		c := ConnectionAdapter{conn: conn}

		require.NoError(t, c.UpgradeTLS(tlsConfig))
		require.True(t, c.TLS)
	}()

	conn, err := net.Dial("tcp", ":1234")
	require.NoError(t, err)
	require.NotNil(t, conn)

	_, err = conn.Write([]byte("lol"))
	require.NoError(t, err)

	wg.Wait()
}

func Test_ConnectionTimeout(t *testing.T) {
	ctx := context.Background()

	sv, err := net.Listen("tcp", ":1234")
	require.NoError(t, err, "should initialize tls server")
	defer sv.Close()

	timeout := 10 * time.Millisecond

	var wg sync.WaitGroup

	go func() {
		defer wg.Done()
		wg.Add(1)
		conn, err := sv.Accept()
		require.NoError(t, err)

		NewConnection(ctx, conn, timeout, nil)

		time.Sleep(timeout)

		one := make([]byte, 1)
		_, err = conn.Read(one)
		require.Error(t, err, "timed-out connection should give error")
	}()

	conn, err := net.Dial("tcp", ":1234")
	require.NoError(t, err)
	require.NotNil(t, conn)

	wg.Wait()
}

func Test_ConnectionReadLine(t *testing.T) {
	ctx := context.Background()

	sv, err := net.Listen("tcp", ":1234")
	require.NoError(t, err, "should initialize tls server")
	defer sv.Close()

	expLine := []byte("test lalala\r\n")
	buffer := []byte("sa6d4as\n6d54a32")

	var wg sync.WaitGroup

	go func() {
		defer wg.Done()
		wg.Add(1)

		conn, err := sv.Accept()
		require.NoError(t, err)

		c := NewConnection(ctx, conn, time.Hour, nil)
		line, err := c.ReadLine()
		require.NoError(t, err)
		require.Equal(t, expLine[:len(expLine)-2], line)
	}()

	conn, err := net.Dial("tcp", ":1234")
	require.NoError(t, err)
	require.NotNil(t, conn)

	_, err = conn.Write(expLine)
	require.NoError(t, err)

	_, err = conn.Write(buffer)
	require.NoError(t, err)

	wg.Wait()
}

func Test_ConnectionReadBytes(t *testing.T) {
	ctx := context.Background()

	sv, err := net.Listen("tcp", ":1234")
	require.NoError(t, err, "should initialize tls server")
	defer sv.Close()

	buffer := []byte("one\ttwo")

	var wg sync.WaitGroup

	go func() {
		defer wg.Done()
		wg.Add(1)

		conn, err := sv.Accept()
		require.NoError(t, err)

		c := NewConnection(ctx, conn, time.Hour, nil)
		line, err := c.ReadBytes('\t')
		require.NoError(t, err)
		require.Equal(t, []byte("one"), line)
	}()

	conn, err := net.Dial("tcp", ":1234")
	require.NoError(t, err)
	require.NotNil(t, conn)

	_, err = conn.Write(buffer)
	require.NoError(t, err)

	wg.Wait()
}
