package interfaces

import (
	"context"

	data "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/data"
)

type Application interface {
	CreateEvent(ctx context.Context, event *data.Event) (*data.Event, error)
	ReadEvent(ctx context.Context, id int) (*data.Event, error)
	UpdateEvent(ctx context.Context, event *data.Event) (*data.Event, error)
	DeleteEvent(ctx context.Context, event *data.Event) (*data.Event, error)
	ListEvents(ctx context.Context) ([]data.Event, error)
	ListNotSheduledEvents(ctx context.Context) ([]data.Event, error)
}
