package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
	api "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/api"
	middleware "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/middleware"
)

var ErrServerRunning = errors.New("HTTP сервер уже запущен")

type HTTPServer struct {
	mu                sync.Mutex
	httpSrv           *http.Server
	app               interfaces.Application
	logger            interfaces.Logger
	host              string
	port              uint16
	readTimeout       time.Duration
	readHeaderTimeout time.Duration
	writeTimeout      time.Duration
	maxHeaderBytes    int
}

func NewHTTPServer(app interfaces.Application,
	host string,
	port uint16,
	readTimeout time.Duration,
	readHeaderTimeout time.Duration,
	writeTimeout time.Duration,
	maxHeaderBytes int,
	logger interfaces.Logger,
) *HTTPServer {
	return &HTTPServer{
		app:               app,
		logger:            logger,
		host:              host,
		port:              port,
		readTimeout:       readTimeout,
		readHeaderTimeout: readHeaderTimeout,
		writeTimeout:      writeTimeout,
		maxHeaderBytes:    maxHeaderBytes,
	}
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello, World!")
}

func (s *HTTPServer) getServer() *http.Server {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.httpSrv
}

func (s *HTTPServer) Start(ctx context.Context) error {
	if ctx == nil {
		return errors.New("контекст запуска HTTP сервера не задан")
	}

	s.mu.Lock()
	mux := http.NewServeMux()
	mux.Handle("/hello", middleware.Instance().Listen(http.HandlerFunc(helloWorld)))
	mux.Handle("/api/", api.Handlers(s.logger, s.app))
	httpSrv := &http.Server{
		Addr:              net.JoinHostPort(s.host, fmt.Sprint(s.port)),
		Handler:           mux,
		ReadTimeout:       s.readTimeout,
		ReadHeaderTimeout: s.readHeaderTimeout,
		WriteTimeout:      s.writeTimeout,
		MaxHeaderBytes:    s.maxHeaderBytes,
		ErrorLog:          log.New(s.logger, "[ERROR]", log.LstdFlags|log.Lmsgprefix),
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
	if s.httpSrv != nil {
		s.mu.Unlock()
		return ErrServerRunning
	}
	s.httpSrv = httpSrv
	s.mu.Unlock()

	// Очищаем поле, когда Serve завершится.
	defer func() {
		s.mu.Lock()
		// Проверка нужна, чтобы случайно не очистить ссылку
		// на другой экземпляр сервера.
		if s.httpSrv == httpSrv {
			s.httpSrv = nil
		}
		s.mu.Unlock()
	}()

	s.logger.Info("Запуск HTTP сервера")

	if err := httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("Ошибка HTTP сервера: %v", err)
		return fmt.Errorf("ошибка HTTP сервера: %w", err)
	}

	return nil
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	if ctx == nil {
		return errors.New("контекст запуска HTTP сервера не задан")
	}
	server := s.getServer()
	if server == nil {
		s.logger.Error("Ошибка при останове HTTP сервера: HTTP сервер не запущен")
		return nil
	}
	s.logger.Info("Останов HTTP сервера")
	err := server.Shutdown(ctx)
	if err != nil {
		s.logger.Error("Ошибка при останове HTTP сервера: %v", err)
		return fmt.Errorf("ошибка при останове HTTP сервера: %w", err)
	}
	return nil
}
