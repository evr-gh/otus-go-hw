package config

import (
	"time"
)

type HTTPConfig struct {
	Host              string
	Port              uint16
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	MaxHeaderBytes    int
}

type HTTPClient struct {
	Host string
	Port uint16
}
