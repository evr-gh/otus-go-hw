package app

import (
	"context"

	data "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/data"
	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
)

type App struct {
	logger  interfaces.Logger
	storage interfaces.Storage
}

func New(logger interfaces.Logger, storage interfaces.Storage) *App {
	return &App{logger, storage}
}

func (a *App) CreateEvent(ctx context.Context, event *data.Event) (*data.Event, error) {
	err := a.storage.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer a.storage.Close(ctx)
	// TODO: in args `ctx context.Context`
	return a.storage.CreateEvent(ctx, event)
}

func (a *App) ReadEvent(ctx context.Context, id int) (*data.Event, error) {
	err := a.storage.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer a.storage.Close(ctx)
	return a.storage.ReadEvent(ctx, id)
}

func (a *App) UpdateEvent(ctx context.Context, e *data.Event) (*data.Event, error) {
	err := a.storage.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer a.storage.Close(ctx)
	return a.storage.UpdateEvent(ctx, e)
}

func (a *App) DeleteEvent(ctx context.Context, e *data.Event) (*data.Event, error) {
	err := a.storage.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer a.storage.Close(ctx)
	return a.storage.DeleteEvent(ctx, e)
}

func (a *App) ListEvents(ctx context.Context) ([]data.Event, error) {
	err := a.storage.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer a.storage.Close(ctx)
	return a.storage.ListEvents(ctx)
}

func (a *App) ListNotSheduledEvents(ctx context.Context) ([]data.Event, error) {
	err := a.storage.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer a.storage.Close(ctx)
	return a.storage.ListNotSheduledEvents(ctx)
}
