package main

import (
	"errors"
	"fmt"
	"os"

	config "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/spf13/viper"
)

var ErrNoConfFile = errors.New("configuration file isn't set (--config <Path to configuration file>)")

func readConfig(configFile string) (*config.CalendarConfig, error) {
	if configFile == "" {
		return nil, ErrNoConfFile
	}

	viper.SetConfigType("yaml")

	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("fail to open configuration file: %w", err)
	}

	viper.ReadConfig(file)
	cmdConfig := config.NewCalendarConfig()
	err = viper.Unmarshal(cmdConfig)
	if err != nil {
		return nil, fmt.Errorf("fail to decode configuration file (%s): %w", configFile, err)
	}
	return cmdConfig, nil
}
