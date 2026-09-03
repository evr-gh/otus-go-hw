package interfaces

import (
	"context"
	"errors"

	data "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/data"
)

var (
	ErrNoData      = errors.New("нет в БД")
	ErrOpInterrupt = errors.New("операция прервана")
	ErrNoEvent     = errors.New("не передана информация по событию")
)

type Storage interface {
	Connect(ctx context.Context) error
	Close() error
	CreateEvent(ctx context.Context, event *data.Event) (*data.Event, error)
	ReadEvent(ctx context.Context, eventID int) (*data.Event, error)
	UpdateEvent(ctx context.Context, event *data.Event) (*data.Event, error)
	DeleteEvent(ctx context.Context, event *data.Event) (*data.Event, error)
	ListEvents(ctx context.Context) ([]data.Event, error)
	ListNotSheduledEvents(ctx context.Context) ([]data.Event, error)
}
