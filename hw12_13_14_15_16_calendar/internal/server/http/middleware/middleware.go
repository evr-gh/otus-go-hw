package middleware

import (
	"net/http"
	"sync"
	"time"

	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
)

type Middleware struct {
	logger interfaces.Logger
}

var (
	middleware *Middleware
	once       sync.Once
)

func Instance() *Middleware {
	if middleware == nil {
		panic("Промежуточное ПО не инициализировано")
	}
	return middleware
}

func Init(logger interfaces.Logger) *Middleware {
	once.Do(func() {
		middleware = &Middleware{}
	})
	middleware.logger = logger

	return middleware
}

func (m Middleware) Listen(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(w)
		StartAt := time.Now()
		handler.ServeHTTP(lrw, r)
		duration := time.Since(StartAt)
		m.logger.Info("Выполнение метода: method=%s[%s]:%s from=%s time=%s code=%v duration=%s",
			r.Method, r.Proto, r.URL.Path, r.RemoteAddr, StartAt, lrw.StatusCode, duration)
	})
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func NewLoggingResponseWriter(writer http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{writer, 0}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
