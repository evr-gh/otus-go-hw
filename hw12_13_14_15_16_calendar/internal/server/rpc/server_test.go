package rpcserver

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	app "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	logger "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	client "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/rpc/client"
	pb "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/rpc/rpcapi"
	storage "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	testServerAddress1 = "localhost:5001"
	testServerAddress2 = "localhost:5002"
)

func newEvent(title string, shd bool) *pb.Event {
	return &pb.Event{
		Title:          title,
		Description:    "integration test event",
		Time:           timestamppb.New(time.Date(2026, time.September, 10, 15, 30, 0, 0, time.UTC)),
		Duration:       durationpb.New(time.Hour),
		Owner:          "rpc-test-user",
		Notifyleadtime: durationpb.New(25 * time.Minute),
		Sheduled:       shd,
	}
}

func assertEventsEqual(t *testing.T, want, got *pb.Event) {
	t.Helper()

	require.NotNil(t, got)

	require.Equal(t, want.Id, got.Id)

	require.Equal(t, want.Title, got.Title)

	require.Equal(t, want.Description, got.Description)

	require.Equal(t, want.Time.AsTime(), got.Time.AsTime())

	require.Equal(t, want.Duration.AsDuration(), got.Duration.AsDuration())

	require.Equal(t, want.Owner, got.Owner)

	require.Equal(t, want.Notifyleadtime.AsDuration(), got.Notifyleadtime.AsDuration())

	require.Equal(t, want.Sheduled, got.Sheduled)
}

func TestRpcServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var once sync.Once
	defer once.Do(cancel)

	logg := logger.New(logger.INFO, os.Stdout)
	storage := storage.New(storage.MemoryStorage, "")
	calendar := app.New(logg, storage)
	grpcServer := NewRPCServer(calendar, logg)

	wg := sync.WaitGroup{}
	wg.Go(func() {
		<-ctx.Done()
		logg.Info("Завершение работы сервера")
		grpcServer.GracefulStop()
	})
	wg.Go(func() {
		logg.Info("Запуск сервера (%v)", testServerAddress1)
		err := grpcServer.Start(ctx, testServerAddress1)
		require.NoError(t, err)
	})

	time.Sleep(5 * time.Second)

	grpcClient := client.Client{}
	logg.Info("Подключение клиента к серверу (%s)", testServerAddress1)
	err := grpcClient.Connect(testServerAddress1)
	require.NoError(t, err)

	pbEvent1 := newEvent("Title 1", false)
	createdEvent1, err := grpcClient.CreateEvent(ctx, pbEvent1)
	require.NoError(t, err)
	pbEvent1.Id = 1
	assertEventsEqual(t, pbEvent1, createdEvent1)

	pbEvent2 := newEvent("Title 2", true)
	createdEvent2, err := grpcClient.CreateEvent(ctx, pbEvent2)
	require.NoError(t, err)
	require.Equal(t, int64(2), createdEvent2.Id)

	pbEvent3 := newEvent("Title 3", false)
	createdEvent3, err := grpcClient.CreateEvent(ctx, pbEvent3)
	require.NoError(t, err)
	require.Equal(t, int64(3), createdEvent3.Id)

	deletedEvent2, err := grpcClient.DeleteEvent(ctx, createdEvent2)
	require.NoError(t, err)
	createdEvent2.Id = 0
	assertEventsEqual(t, createdEvent2, deletedEvent2)

	ident := pb.Id{}
	ident.Id = 3
	pbEvent3Copy, err := grpcClient.ReadEvent(ctx, &ident)
	require.NoError(t, err)
	assertEventsEqual(t, createdEvent3, pbEvent3Copy)

	createdEvent3.Sheduled = true
	pbEvent3Updated, err := grpcClient.UpdateEvent(ctx, createdEvent3)
	require.NoError(t, err)
	assertEventsEqual(t, createdEvent3, pbEvent3Updated)

	events, err := grpcClient.ListEvents(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(events))

	count := 0
	for _, e := range events {
		switch e.Id {
		case 1:
			assertEventsEqual(t, createdEvent1, e)
			count++
		case 3:
			assertEventsEqual(t, createdEvent3, e)
			count++
		}
	}

	require.Equal(t, 2, count)

	events, err = grpcClient.ListNotSheduledEvents(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(events))
	assertEventsEqual(t, createdEvent1, events[0])

	grpcClient.Close()

	once.Do(cancel)
	wg.Wait()
}

func TestInterceptorLogging(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var once sync.Once
	defer once.Do(cancel)

	outputInto := &bytes.Buffer{}
	logg := logger.New(logger.DEBUG, outputInto)

	storage := storage.New(storage.MemoryStorage, "")
	calendar := app.New(logg, storage)
	grpcServer := NewRPCServer(calendar, logg)

	wg := sync.WaitGroup{}
	wg.Go(func() {
		<-ctx.Done()
		logg.Info("Завершение работы сервера")
		grpcServer.GracefulStop()
	})
	wg.Go(func() {
		logg.Info("Запуск сервера (%v)", testServerAddress2)
		err := grpcServer.Start(ctx, testServerAddress2)
		require.NoError(t, err)
	})
	//
	grpcClient := client.Client{}
	grpcClient.Connect(testServerAddress2)

	pbEvent1 := newEvent("Title 1", false)
	createdEvent1, err := grpcClient.CreateEvent(ctx, pbEvent1)
	require.NoError(t, err)
	require.Equal(t, int64(1), createdEvent1.Id)

	pbEvent2 := newEvent("Title 2", false)
	createdEvent2, err := grpcClient.CreateEvent(ctx, pbEvent2)
	require.NoError(t, err)
	require.Equal(t, int64(2), createdEvent2.Id)

	_, err = grpcClient.ListEvents(ctx)
	require.NoError(t, err)

	outputted := outputInto.String()

	fmt.Println(outputted)
	require.True(t, strings.Contains(outputted, "[INFO] Запуск gRPC сервера: address=\"localhost:5002\""))
	require.True(t, strings.Contains(outputted, "[INFO] Выполнение метода: method=\"/calendar.Application/CreateEvent\" out={id:1  title:\"Title 1\""))
	require.True(t, strings.Contains(outputted, "[INFO] Выполнение метода: method=\"/calendar.Application/CreateEvent\" out={id:2  title:\"Title 2\""))
	require.True(t, strings.Contains(outputted, "[INFO] Начало gRPC потока: method=\"/calendar.Application/ListEvents\" client_stream=false server_stream=true"))
	require.True(t, strings.Contains(outputted, "[INFO] Конец gRPC потока: method=\"/calendar.Application/ListEvents\" code=OK"))

	grpcClient.Close()
	once.Do(cancel)
	wg.Wait()
}
