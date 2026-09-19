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
	"github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/middleware"
	storage "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/spf13/pflag"
)

var configFile string

func init() {
	pflag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
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
		os.Exit(1)
	}

	var logFile *os.File
	var logg *logger.Logger

	if cmdConfig.Logger.File != "" {
		logFile, err = os.Create(cmdConfig.Logger.File)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Не удалось создать файл для сохранения лога: %v", err)
			os.Exit(1)
		}
		defer logFile.Close()
		logg = logger.New(cmdConfig.Logger.Level, logFile)
	} else {
		logg = logger.New(cmdConfig.Logger.Level, os.Stdout)
	}

	stg, err := storage.New(cmdConfig.Storage.Type, cmdConfig.Storage.DSN)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if logFile != nil {
			logFile.Close()
		}
		os.Exit(1) //nolint:gocritic
	}

	calendar := app.New(logg, stg)

	middleware.Init(logg)
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

	err = calendar.Start(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		stop()
		if logFile != nil {
			logFile.Close()
		}
		os.Exit(1)
	}
	defer calendar.Close()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("Не удалось остановить HTTP сервер: %v", err.Error())
		}
	}()

	logg.Info("Начало работы сервиса \"Календарь\"")

	if err := server.Start(ctx); err != nil {
		logg.Error("Не удалось запустить HTTP сервер: %v", err.Error())
		calendar.Close()
		stop()
		if logFile != nil {
			logFile.Close()
		}
		os.Exit(1)
	}
}
