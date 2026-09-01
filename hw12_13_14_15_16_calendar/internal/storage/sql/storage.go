package sqlstorage

import (
	"context"

	data "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/data"
)

type Storage struct { // TODO
}

func New(dsn string) *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event *data.Event) (*data.Event, error) {
	return event, nil
}

func (s *Storage) ReadEvent(ctx context.Context, eventID int) (*data.Event, error) {
	return nil, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event *data.Event) (*data.Event, error) {
	return event, nil
}

func (s *Storage) DeleteEvent(ctx context.Context, event *data.Event) (*data.Event, error) {
	return event, nil
}

func (s *Storage) ListEvents(ctx context.Context) ([]data.Event, error) {
	return []data.Event{}, nil
}

func (s *Storage) ListNotSheduledEvents(ctx context.Context) ([]data.Event, error) {
	return []data.Event{}, nil
}
