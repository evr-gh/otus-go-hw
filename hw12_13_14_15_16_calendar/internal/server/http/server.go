package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
	middleware "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/middleware"
)

type Server struct {
	httpSrv *http.Server
	app     interfaces.Application
	logger  interfaces.Logger
	context context.Context
}

func NewServer(app interfaces.Application,
	host string,
	port uint16,
	readTimeout time.Duration,
	readHeaderTimeout time.Duration,
	writeTimeout time.Duration,
	maxHeaderBytes int,
	logger interfaces.Logger,
) *Server {
	mux := http.NewServeMux()
	mux.Handle("/hello", middleware.Instance().Listen(http.HandlerFunc(helloWord)))
	server := new(Server)
	server.app = app
	server.logger = logger
	server.context = context.Background()
	server.httpSrv = &http.Server{
		Addr:              net.JoinHostPort(host, fmt.Sprint(port)),
		Handler:           mux,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
		ErrorLog:          log.New(logger, "[ERROR]", log.LstdFlags|log.Lmsgprefix),
		BaseContext: func(_ net.Listener) context.Context {
			return server.context
		},
	}
	return server
}

func helloWord(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello, World!")
}

func (s *Server) Start(ctx context.Context) error {
	s.context = ctx
	s.logger.Info("Запуск HTTP сервера")

	if err := s.httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("Ошибка HTTP сервера: %v", err)
		return fmt.Errorf("ошибка HTTP сервера: %w", err)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Останов HTTP сервера")
	err := s.httpSrv.Shutdown(ctx)
	if err != nil {
		s.logger.Error("Ошибка при останове HTTP сервера: %v", err)
		return fmt.Errorf("ошибка при останове HTTP сервера: %w", err)
	}
	return nil
}
