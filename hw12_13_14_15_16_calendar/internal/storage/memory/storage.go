package memorystorage

import (
	"context"
	"fmt"
	"sync"

	models "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/data"
	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
)

type inMemoryDatabase map[int]*models.Event

type Storage struct {
	data    inMemoryDatabase
	mu      sync.RWMutex
	counter int
}

func New() *Storage {
	return &Storage{make(inMemoryDatabase), sync.RWMutex{}, 0}
}

func (s *Storage) Connect(_ context.Context) error {
	return nil
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) CreateEvent(_ context.Context, event *models.Event) (*models.Event, error) {
	if event == nil {
		return event, fmt.Errorf("не создано событие: %w", interfaces.ErrNoEvent)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	event.ID = s.counter
	event.Sheduled = false
	s.data[event.ID] = event
	return event, nil
}

func (s *Storage) ReadEvent(_ context.Context, eventID int) (*models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, exists := s.data[eventID]
	if exists {
		return event, nil
	}
	return nil, fmt.Errorf("не получено событие с ID=%v: %w", eventID, interfaces.ErrNoData)
}

func (s *Storage) UpdateEvent(_ context.Context, event *models.Event) (*models.Event, error) {
	if event == nil {
		return event, fmt.Errorf("не обновлено событие: %w", interfaces.ErrNoEvent)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.data[event.ID]
	if exists {
		s.data[event.ID] = event
		return event, nil
	}
	return event, fmt.Errorf("не обновлено событие с ID=%v: %w", event.ID, interfaces.ErrNoData)
}

func (s *Storage) DeleteEvent(_ context.Context, event *models.Event) (*models.Event, error) {
	if event == nil {
		return event, fmt.Errorf("не удалено событие: %w", interfaces.ErrNoEvent)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.data[event.ID]
	if exists {
		delete(s.data, event.ID)
		event.ID = 0
		return event, nil
	}
	return event, nil
}

func (s *Storage) ListEvents(cntx context.Context) ([]models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]models.Event, 0, len(s.data))
	for _, event := range s.data {
		select {
		case <-cntx.Done():
			return events, fmt.Errorf("список событий не получен: %w", interfaces.ErrOpInterrupt)
		default:
			events = append(events, *event)
		}
	}
	return events, nil
}

func (s *Storage) ListNotSheduledEvents(cntx context.Context) ([]models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]models.Event, 0, len(s.data)/4)
	for _, event := range s.data {
		select {
		case <-cntx.Done():
			return events, fmt.Errorf("список незапланированных событий не получен: %w", interfaces.ErrOpInterrupt)
		default:
			if !event.Sheduled {
				events = append(events, *event)
			}
		}
	}
	return events, nil
}
