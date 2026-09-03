package main

import (
	"bytes"
	"context"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		lc := &net.ListenConfig{}
		l, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
			defer stop()

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect(ctx))
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send(ctx)
			require.NoError(t, err)

			err = client.Receive(ctx)
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
}

func TestAdd1TelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		lc := &net.ListenConfig{}
		l, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
			defer stop()

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect(ctx))
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send(ctx)
			require.NoError(t, err)

			err = client.Receive(ctx)
			require.NoError(t, err)
			require.Equal(t, "", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			require.NoError(t, conn.Close())
		}()

		wg.Wait()
	})
}

func TestAdd2TelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
		defer stop()

		timeout, err := time.ParseDuration("1s")
		require.NoError(t, err)

		client := NewTelnetClient("127.0.0.1:1234", timeout, io.NopCloser(in), out)
		err = client.Connect(ctx)
		require.Error(t, err)
		require.Equal(t, err.Error(), "dial tcp 127.0.0.1:1234: connect: connection refused")

		defer func() { require.NoError(t, client.Close()) }()
	})
}
