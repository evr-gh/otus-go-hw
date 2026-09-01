package memorystorage

import (
	"context"
	"fmt"
	"sync"

	data "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/data"
	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
)

type inMemoryDatabase map[int]*data.Event

type Storage struct {
	data    inMemoryDatabase
	mu      sync.RWMutex
	counter int
}

func New() *Storage {
	return &Storage{make(inMemoryDatabase), sync.RWMutex{}, 0}
}

func (s *Storage) Connect(ctx context.Context) error {
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event *data.Event) (*data.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	event.ID = s.counter
	event.Sheduled = false
	s.data[event.ID] = event
	return event, nil
}

func (s *Storage) ReadEvent(ctx context.Context, eventID int) (*data.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, exists := s.data[eventID]
	if exists {
		return event, nil
	}
	return nil, fmt.Errorf("не получено событие с ID=%v: %w", eventID, interfaces.ErrNoData)
}

func (s *Storage) UpdateEvent(ctx context.Context, event *data.Event) (*data.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.data[event.ID]
	if exists {
		s.data[event.ID] = event
		return event, nil
	}
	return event, fmt.Errorf("не обновлено событие с ID=%v: %w", event.ID, interfaces.ErrNoData)
}

func (s *Storage) DeleteEvent(ctx context.Context, event *data.Event) (*data.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.data[event.ID]
	if exists {
		delete(s.data, event.ID)
		return event, nil
	}
	return event, fmt.Errorf("не удалено событие с ID=%v: %w", event.ID, interfaces.ErrNoData)
}

func (s *Storage) ListEvents(ctx context.Context) ([]data.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]data.Event, 0, len(s.data))
	for _, event := range s.data {
		events = append(events, *event)
	}
	return events, nil
}

func (s *Storage) ListNotSheduledEvents(ctx context.Context) ([]data.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]data.Event, 0, len(s.data)/4)
	for _, event := range s.data {
		if !event.Sheduled {
			events = append(events, *event)
		}
	}
	return events, nil
}
