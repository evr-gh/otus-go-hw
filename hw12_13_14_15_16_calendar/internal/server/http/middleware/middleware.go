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
		panic("Middleware was not init by `Init(logger interfaces.Logger)`.")
	}
	return middleware
}

func Init(logger interfaces.Logger) *Middleware {
	once.Do(func() {
		middleware = &Middleware{}
		middleware.logger = logger
	})
	return middleware
}

func (m Middleware) Listen(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		StartAt := time.Now()
		lrw := NewLoggingResponseWriter(w)
		handler.ServeHTTP(lrw, r)
		a := struct {
			ClientIPAddress string
			StartAt         time.Time
			HTTPMethod      string
			HTTPVersion     string
			URLPath         string
			StatusCode      int
			Latency         time.Duration
		}{
			ClientIPAddress: r.RemoteAddr,
			StartAt:         StartAt,
			HTTPMethod:      r.Method,
			HTTPVersion:     r.Proto,
			URLPath:         r.URL.Path,
			StatusCode:      lrw.StatusCode,
			Latency:         time.Since(StartAt),
		}
		m.logger.Info("%+v", a)
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
