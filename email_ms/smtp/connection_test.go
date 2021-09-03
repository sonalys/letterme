package smtp

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
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

func Test_ConnectionReadEnvelope(t *testing.T) {
	ctx := context.Background()
	sv, err := net.Listen("tcp", ":1234")
	require.NoError(t, err, "should initialize tls server")
	defer sv.Close()

	type testCase struct {
		name           string
		dataToSend     []byte
		expectedBuffer []byte
		expectedError  error
	}

	testList := []testCase{
		{
			name:           "envelope too big",
			dataToSend:     make([]byte, maxEnvelopeDataSize+1),
			expectedBuffer: nil,
			expectedError:  errors.New("too big"),
		},
		{
			name:           "all ok",
			dataToSend:     []byte("hello world\r\n.\r\n"),
			expectedBuffer: []byte("hello world\r\n.\r\n"),
			expectedError:  nil,
		},
	}

	for _, tC := range testList {
		t.Run(tC.name, func(t *testing.T) {
			var wg sync.WaitGroup

			go func() {
				defer wg.Done()
				wg.Add(1)

				conn, err := sv.Accept()
				require.NoError(t, err)
				defer conn.Close()

				c := NewConnection(ctx, conn, time.Hour, nil)

				out, err := c.ReadEnvelope()
				if tC.expectedError == nil {
					require.NoError(t, err)
					require.NotNil(t, out)
				} else {
					require.EqualError(t, err, tC.expectedError.Error())
					return
				}

				gotBytes, err := io.ReadAll(out)
				require.NoError(t, err)

				if tC.expectedBuffer != nil {
					require.Equal(t, tC.expectedBuffer, gotBytes)
				}
			}()

			conn, err := net.Dial("tcp", ":1234")
			require.NoError(t, err)
			require.NotNil(t, conn)

			size := len(tC.dataToSend)
			step := 4096

			for startOffset := 0; startOffset < size; startOffset += step {
				endOffset := startOffset + step
				if endOffset > size {
					endOffset = size
				}

				_, err := conn.Write(tC.dataToSend[startOffset : startOffset+endOffset])
				if err != nil {
					break
				}
			}

			wg.Wait()
		})
	}
}

func Test_ConnectionAddBuffer(t *testing.T) {
	ctx := context.Background()
	sv, err := net.Listen("tcp", ":1234")
	require.NoError(t, err, "should initialize tls server")
	defer sv.Close()

	go func() {
		conn, err := sv.Accept()
		require.NoError(t, err)
		require.NotNil(t, conn)

		c := NewConnection(ctx, conn, time.Hour, nil)

		c.AddBuffer([]byte("world"))
		c.AddBuffer("hello", []byte("world"))

		err = c.Flush()
		require.NoError(t, err)

		c.Close()
	}()

	conn, err := net.Dial("tcp", ":1234")
	require.NoError(t, err)
	require.NotNil(t, conn)

	got := make([]byte, 0, 64)

	for {
		buffer := make([]byte, 4096)
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			break
		}

		got = append(got, buffer[:bytesRead]...)
	}

	expected := []byte("helloworld\r\n")
	require.Equal(t, expected, got)

}
