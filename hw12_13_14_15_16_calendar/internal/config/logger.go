package config

import "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"

type LoggerConfig struct {
	Level interfaces.LogLevel
	File  string
}
