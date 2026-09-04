package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	app "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	logger "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http"
	storage "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/spf13/pflag"
)

var configFile string

func init() {
	pflag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	pflag.Parse()

	if pflag.Arg(0) == "version" {
		printVersion()
		return
	}

	cmdConfig, err := readConfig(configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	logg := logger.New(cmdConfig.Logger.Level, os.Stdout)

	storage := storage.New(cmdConfig.Storage.Type, cmdConfig.Storage.DSN)
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(calendar,
		cmdConfig.HTTP.Host,
		cmdConfig.HTTP.Port,
		cmdConfig.HTTP.ReadTimeout,
		cmdConfig.HTTP.ReadHeaderTimeout,
		cmdConfig.HTTP.WriteTimeout,
		cmdConfig.HTTP.MaxHeaderBytes,
		logg)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: %v", err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: %v", err.Error())
		stop()
		os.Exit(1) //nolint:gocritic
	}
}
