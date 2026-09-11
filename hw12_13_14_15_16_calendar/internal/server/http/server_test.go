package httpserver

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	logger "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	middleware "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/middleware"
	"github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

const (
	host string = "localhost"
)

func TestServerCode(t *testing.T) {
	const port uint16 = 5021

	outputInto := &bytes.Buffer{}
	logg := logger.New(logger.INFO, outputInto)
	middleware.Init(logg)

	httpServer := NewHTTPServer(nil, host, port, 10*time.Second, 11*time.Second, 12*time.Second, 65536, logg)
	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Go(func() {
		if err := httpServer.Start(ctx); err != nil {
			logg.Error("%v", err.Error())
		}
	})

	time.Sleep(5 * time.Second)

	url := fmt.Sprintf("http://%s:%d/hello", host, port)
	client := &http.Client{}
	request, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	resp, err := client.Do(request)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, 200, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, "Hello, World!\n", string(body))

	timeoutCtx, cancelByTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelByTimeout()
	err = httpServer.Stop(timeoutCtx)
	require.NoError(t, err)

	outputted := outputInto.String()

	require.Truef(t, strings.Contains(outputted, "[INFO] Запуск HTTP сервера"), outputted)
	require.Truef(t, strings.Contains(outputted, "[INFO] Выполнение метода:"+
		" method=POST[HTTP/1.1]:/hello from=127.0.0.1"), outputted)
	require.Truef(t, strings.Contains(outputted, "Останов HTTP сервера"), outputted)
}

func TestServerErrCode(t *testing.T) {
	const port uint16 = 5022
	outputInto := &bytes.Buffer{}
	logg := logger.New(logger.INFO, outputInto)
	middleware.Init(logg)

	httpServer := NewHTTPServer(nil, host, port, 10*time.Second, 11*time.Second, 12*time.Second, 65536, logg)
	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Go(func() {
		if err := httpServer.Start(ctx); err != nil {
			logg.Error("%v", err.Error())
		}
	})

	time.Sleep(2 * time.Second)

	url := fmt.Sprintf("http://%s:%d/err", host, port)
	client := &http.Client{}
	request, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	resp, err := client.Do(request)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, 404, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, "404 page not found\n", string(body))

	timeoutCtx, cancelByTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelByTimeout()
	err = httpServer.Stop(timeoutCtx)
	require.NoError(t, err)

	outputted := outputInto.String()
	require.True(t, strings.Contains(outputted, "Запуск HTTP сервера"))
	require.True(t, strings.Contains(outputted, "Останов HTTP сервера"))
}

func TestServerStopNotStarted(t *testing.T) {
	const port uint16 = 5023
	outputInto := &bytes.Buffer{}
	logg := logger.New(logger.INFO, outputInto)
	middleware.Init(logg)

	httpServer := NewHTTPServer(nil, host, port, 10*time.Second, 11*time.Second, 12*time.Second, 65536, logg)

	err := httpServer.Stop(context.Background())
	require.NoError(t, err)
}

func TestServerStopNormally(t *testing.T) {
	const port uint16 = 5024
	outputInto := &bytes.Buffer{}
	logg := logger.New(logger.INFO, outputInto)
	middleware.Init(logg)

	httpServer := NewHTTPServer(nil, host, port, 10*time.Second, 11*time.Second, 12*time.Second, 65536, logg)

	wg := sync.WaitGroup{}
	wg.Go(func() {
		time.Sleep(3 * time.Second)
		err := httpServer.Stop(context.Background())
		require.NoError(t, err)
	})
	err := httpServer.Start(context.Background())
	require.NoError(t, err)
	wg.Wait()
}

func TestServerStopBySignal(t *testing.T) {
	const port uint16 = 5025
	ctx, ctxCancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT)
	defer ctxCancel()

	outputInto := &bytes.Buffer{}
	logg := logger.New(logger.INFO, outputInto)
	middleware.Init(logg)

	httpServer := NewHTTPServer(nil, host, port, 10*time.Second, 11*time.Second, 12*time.Second, 65536, logg)

	wg := sync.WaitGroup{}
	wg.Go(func() {
		err := httpServer.Start(ctx)
		require.NoError(t, err)
	})
	wg.Go(func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()
		err := httpServer.Stop(shutdownCtx)
		require.NoError(t, err)
	})
	pid, _, _ := syscall.Syscall(syscall.SYS_GETPID, 0, 0, 0)
	process, _ := os.FindProcess(int(pid))
	process.Signal(syscall.SIGHUP)
	wg.Wait()
}

func TestServerStopBySignalAfterDelay(t *testing.T) {
	const port uint16 = 5026
	ctx, ctxCancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT)
	defer ctxCancel()

	outputInto := &bytes.Buffer{}
	logg := logger.New(logger.INFO, outputInto)
	middleware.Init(logg)

	httpServer := NewHTTPServer(nil, host, port, 10*time.Second, 11*time.Second, 12*time.Second, 65536, logg)

	wg := sync.WaitGroup{}
	wg.Go(func() {
		err := httpServer.Start(ctx)
		require.NoError(t, err)
	})
	wg.Go(func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()
		err := httpServer.Stop(shutdownCtx)
		require.NoError(t, err)
	})
	time.Sleep(3 * time.Second)
	pid, _, _ := syscall.Syscall(syscall.SYS_GETPID, 0, 0, 0)
	process, _ := os.FindProcess(int(pid))
	process.Signal(syscall.SIGHUP)
	wg.Wait()
}

func TestServerStopByCancel(t *testing.T) {
	const port uint16 = 5027
	ctx, ctxCancel := context.WithCancel(context.Background())
	outputInto := &bytes.Buffer{}

	logg := logger.New(logger.INFO, outputInto)
	middleware.Init(logg)

	calendarApp := app.New(logg, storage.New("memory", ""))
	httpServer := NewHTTPServer(calendarApp, host, port, 10*time.Second, 11*time.Second, 12*time.Second, 65536, logg)
	wg := sync.WaitGroup{}
	wg.Go(func() {
		err := httpServer.Start(ctx)
		require.NoError(t, err)
	})
	wg.Go(func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()
		err := httpServer.Stop(shutdownCtx)
		require.NoError(t, err)
	})

	time.Sleep(3 * time.Second)
	ctxCancel()
	wg.Wait()
}
