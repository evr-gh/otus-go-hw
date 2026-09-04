package interfaces

import (
	"context"

	models "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/data"
)

type Application interface {
	CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error)
	ReadEvent(ctx context.Context, id int) (*models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) (*models.Event, error)
	DeleteEvent(ctx context.Context, event *models.Event) (*models.Event, error)
	ListEvents(ctx context.Context) ([]models.Event, error)
	ListNotSheduledEvents(ctx context.Context) ([]models.Event, error)
}
