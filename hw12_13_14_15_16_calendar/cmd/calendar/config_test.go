package main

import (
	"testing"
	"time"

	logger "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		cmdConfig, err := readConfig("../../configs/calendar.test.yaml")
		require.NoError(t, err)

		require.Equal(t, "0.0.0.0", cmdConfig.HTTP.Host)
		require.Equal(t, uint16(8000), cmdConfig.HTTP.Port)
		timeout, err := time.ParseDuration("5s")
		require.NoError(t, err)
		require.Equal(t, timeout, cmdConfig.HTTP.ReadTimeout)
		timeout, err = time.ParseDuration("6s")
		require.NoError(t, err)
		require.Equal(t, timeout, cmdConfig.HTTP.ReadHeaderTimeout)
		timeout, err = time.ParseDuration("7s")
		require.NoError(t, err)
		require.Equal(t, timeout, cmdConfig.HTTP.WriteTimeout)
		require.Equal(t, 65000, cmdConfig.HTTP.MaxHeaderBytes)

		require.Equal(t, "1.2.3.4", cmdConfig.RPC.Host)
		require.Equal(t, uint16(5000), cmdConfig.RPC.Port)

		require.Equal(t, "pgsqldtb", cmdConfig.Storage.Type)
		require.Equal(t,
			"user=user password=userpasswd host=db_host database=clendar search_path=calendar sslmode=disable port=5432",
			cmdConfig.Storage.DSN)

		require.Equal(t, logger.LogLevel("INFO"), cmdConfig.Logger.Level)
	})

	t.Run("no conf file", func(t *testing.T) {
		cmdConfig, err := readConfig("")
		require.Error(t, err)
		require.Equal(t, err.Error(), "configuration file isn't set (--config <Path to configuration file>)")
		require.Nil(t, cmdConfig)
	})

	t.Run("no exist conf file", func(t *testing.T) {
		cmdConfig, err := readConfig("../../configs/calendar1.yaml")
		require.Error(t, err)
		//nolint:lll
		require.Equal(t, err.Error(), "fail to open configuration file: open ../../configs/calendar1.yaml: no such file or directory")
		require.Nil(t, cmdConfig)
	})
	t.Run("invalid file", func(t *testing.T) {
		cmdConfig, err := readConfig("../../configs/calendar.invalid.yaml")
		require.Error(t, err)
		//nolint:lll
		require.Equal(t, err.Error(), "fail to decode configuration file (../../configs/calendar.invalid.yaml): decoding failed due to the following error(s):\n\n'HTTP.Port' cannot parse value as 'uint16': strconv.ParseUint: invalid syntax")
		require.Nil(t, cmdConfig)
	})
}
