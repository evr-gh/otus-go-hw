package interfaces

import (
	"io"
)

type Logger interface {
	io.Writer
	Debug(msg string, a ...any)
	Info(msg string, a ...any)
	Warning(msg string, a ...any)
	Error(msg string, a ...any)
}
