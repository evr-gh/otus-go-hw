package logger

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type LogLevel string

const (
	DEBUG   LogLevel = "DEBUG"
	INFO    LogLevel = "INFO"
	WARNING LogLevel = "WARNING"
	ERROR   LogLevel = "ERROR"
)

type Logger struct {
	level  LogLevel
	writer io.Writer
}

func New(level LogLevel, writer io.Writer) *Logger {
	return &Logger{level, writer}
}

func (l Logger) saveMsg(level LogLevel, template string, a ...any) {
	var buildedString strings.Builder
	fmt.Fprintf(&buildedString, "%s [%s] ", level, time.Now().UTC().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&buildedString, template, a...)
	if !strings.HasSuffix(template, "\n") {
		buildedString.WriteString("\n")
	}
	l.writer.Write([]byte(buildedString.String()))
}

func (l Logger) Debug(template string, a ...any) {
	if l.level == DEBUG {
		l.saveMsg(DEBUG, template, a...)
	}
}

func (l Logger) Info(template string, a ...any) {
	if l.level == INFO || l.level == DEBUG {
		l.saveMsg(INFO, template, a...)
	}
}

func (l Logger) Warning(template string, a ...any) {
	if l.level == WARNING || l.level == INFO || l.level == DEBUG {
		l.saveMsg(WARNING, template, a...)
	}
}

func (l Logger) Error(template string, a ...any) {
	if l.level == ERROR || l.level == WARNING || l.level == INFO || l.level == DEBUG {
		l.saveMsg(ERROR, template, a...)
	}
}

func (l Logger) Log(template string, a ...any) {
	l.saveMsg(l.level, template, a...)
}
