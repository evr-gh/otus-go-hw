package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
	middleware "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/middleware"
)

type Server struct {
	server *http.Server
	app    interfaces.Application
	logger interfaces.Logger
}

type Logger interface { // TODO
}

type Application interface { // TODO
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
	mux.Handle("/", middleware.Instance().Listen(http.HandlerFunc(hellowWord)))
	server := http.Server{
		Addr:              net.JoinHostPort(host, fmt.Sprint(port)),
		Handler:           mux,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}
	return &Server{&server, app, logger}
}

func hellowWord(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello Word!"))
}

func (s *Server) Start(ctx context.Context) error {
	_ = ctx // TODO: for what?
	s.logger.Info("Запуск HTTP сервера")

	if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("Ошибка HTTP сервера: %v", err)
		return fmt.Errorf("ошибка HTTP сервера: %w", err)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Останов HTTP сервера")
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Error("Ошибка при останове HTTP сервера: %v", err)
		return fmt.Errorf("ошибка при останове HTTP сервера: %w", err)
	}
	return nil
}
