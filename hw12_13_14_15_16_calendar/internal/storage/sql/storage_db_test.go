package sqlstorage

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	models "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/data"
	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
)

const defaultTestDSN = "postgres://user:userpasswd@localhost:5432/calendar?sslmode=disable&search_path=calendar"

func newTestStorage(t *testing.T) *Storage {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		dsn = defaultTestDSN
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	storage := New("pgx", dsn)

	if err := storage.Connect(ctx); err != nil {
		t.Fatalf("не удалось подключиться к тестовой БД: %v", err)
	}

	t.Cleanup(func() {
		if err := storage.Close(); err != nil {
			t.Errorf("не удалось закрыть соединение с тестовой БД: %v", err)
		}
	})

	truncateEvents(t, storage)

	return storage
}

func truncateEvents(t *testing.T, storage *Storage) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := storage.db.ExecContext(
		ctx,
		`TRUNCATE TABLE events RESTART IDENTITY`,
	); err != nil {
		t.Fatalf("не удалось очистить таблицу events: %v", err)
	}
}

func newEvent(title string, sheduled bool) *models.Event {
	return &models.Event{
		Title:          title,
		Description:    "integration test event",
		Time:           time.Date(2026, time.September, 7, 15, 30, 0, 0, time.UTC),
		Duration:       time.Hour,
		Owner:          "integration-test-user",
		NotifyLeadTime: 15 * time.Minute,
		Sheduled:       sheduled,
	}
}

func assertEventsEqual(t *testing.T, want, got *models.Event) {
	t.Helper()

	if got == nil {
		t.Fatal("получено nil-событие")
	}

	if got.ID != want.ID {
		t.Errorf("ID: ожидалось %v, получено %v", want.ID, got.ID)
	}

	if got.Title != want.Title {
		t.Errorf("Title: ожидалось %q, получено %q", want.Title, got.Title)
	}

	if got.Description != want.Description {
		t.Errorf(
			"Description: ожидалось %q, получено %q",
			want.Description,
			got.Description,
		)
	}

	if !got.Time.Equal(want.Time) {
		t.Errorf("Time: ожидалось %v, получено %v", want.Time, got.Time)
	}

	if got.Duration != want.Duration {
		t.Errorf(
			"Duration: ожидалось %v, получено %v",
			want.Duration,
			got.Duration,
		)
	}

	if got.Owner != want.Owner {
		t.Errorf("Owner: ожидалось %v, получено %v", want.Owner, got.Owner)
	}

	if got.NotifyLeadTime != want.NotifyLeadTime {
		t.Errorf(
			"NotifyLeadTime: ожидалось %v, получено %v",
			want.NotifyLeadTime,
			got.NotifyLeadTime,
		)
	}

	if got.Sheduled != want.Sheduled {
		t.Errorf(
			"Sheduled: ожидалось %v, получено %v",
			want.Sheduled,
			got.Sheduled,
		)
	}
}

func TestStorageCreateAndReadEvent(t *testing.T) {
	storage := newTestStorage(t)
	ctx := context.Background()

	event := newEvent("create and read test", false)

	created, err := storage.CreateEvent(ctx, event)
	if err != nil {
		t.Fatalf("CreateEvent вернул ошибку: %v", err)
	}

	if created == nil {
		t.Fatal("CreateEvent вернул nil")
	}

	if created.ID <= 0 {
		t.Fatalf("ожидался положительный ID, получено %d", created.ID)
	}

	read, err := storage.ReadEvent(ctx, created.ID)
	if err != nil {
		t.Fatalf("ReadEvent вернул ошибку: %v", err)
	}

	assertEventsEqual(t, created, read)
}

func TestStorageReadEventNotFound(t *testing.T) {
	storage := newTestStorage(t)

	event, err := storage.ReadEvent(context.Background(), 999999)
	if err == nil {
		t.Fatal("ожидалась ошибка для отсутствующего события")
	}

	if event != nil {
		t.Errorf("ожидалось nil-событие, получено %#v", event)
	}

	if !errors.Is(err, interfaces.ErrNoData) {
		t.Fatalf(
			"ожидалась ошибка ErrNoData, получено: %v",
			err,
		)
	}
}

func TestStorageUpdateEvent(t *testing.T) {
	storage := newTestStorage(t)
	ctx := context.Background()

	created, err := storage.CreateEvent(ctx, newEvent("old title", false))
	if err != nil {
		t.Fatalf("не удалось создать событие: %v", err)
	}

	created.Title = "updated title"
	created.Description = "updated description"
	created.Duration = 30 * time.Minute
	created.NotifyLeadTime = 3 * time.Hour
	created.Sheduled = true
	created.Time = created.Time.Add(5 * time.Minute)

	updated, err := storage.UpdateEvent(ctx, created)
	if err != nil {
		t.Fatalf("UpdateEvent вернул ошибку: %v", err)
	}

	assertEventsEqual(t, created, updated)

	read, err := storage.ReadEvent(ctx, created.ID)
	if err != nil {
		t.Fatalf("не удалось прочитать обновлённое событие: %v", err)
	}

	assertEventsEqual(t, created, read)
}

func TestStorageUpdateEventNil(t *testing.T) {
	storage := newTestStorage(t)

	event, err := storage.UpdateEvent(context.Background(), nil)
	if err == nil {
		t.Fatal("ожидалась ошибка для nil-события")
	}

	if event != nil {
		t.Errorf("ожидалось nil-событие, получено %#v", event)
	}

	if !errors.Is(err, interfaces.ErrNoEvent) {
		t.Fatalf(
			"ожидалась ошибка ErrNoEvent, получено: %v",
			err,
		)
	}
}

func TestStorageUpdateEventNotFound(t *testing.T) {
	storage := newTestStorage(t)

	event := newEvent("missing event", false)
	event.ID = 999999

	updated, err := storage.UpdateEvent(context.Background(), event)
	if err == nil {
		t.Fatal("ожидалась ошибка для отсутствующего события")
	}

	if updated != event {
		t.Error("UpdateEvent должен вернуть переданный объект события")
	}

	if !errors.Is(err, interfaces.ErrNoData) {
		t.Fatalf(
			"ожидалась ошибка ErrNoData, получено: %v",
			err,
		)
	}
}

func TestStorageDeleteEvent(t *testing.T) {
	storage := newTestStorage(t)
	ctx := context.Background()

	created, err := storage.CreateEvent(ctx, newEvent("delete test", false))
	if err != nil {
		t.Fatalf("не удалось создать событие: %v", err)
	}

	eventID := created.ID

	deleted, err := storage.DeleteEvent(ctx, created)
	if err != nil {
		t.Fatalf("DeleteEvent вернул ошибку: %v", err)
	}

	if deleted == nil {
		t.Fatal("DeleteEvent вернул nil")
	}

	if deleted.ID != 0 {
		t.Errorf("после удаления ожидался ID=0, получено %d", deleted.ID)
	}

	read, err := storage.ReadEvent(ctx, eventID)
	if err == nil {
		t.Fatalf("удалённое событие всё ещё доступно: %#v", read)
	}

	if !errors.Is(err, interfaces.ErrNoData) {
		t.Fatalf(
			"после удаления ожидалась ошибка ErrNoData, получено: %v",
			err,
		)
	}
}

func TestStorageDeleteEventValidation(t *testing.T) {
	storage := newTestStorage(t)
	ctx := context.Background()

	t.Run("nil event", func(t *testing.T) {
		event, err := storage.DeleteEvent(ctx, nil)
		if err == nil {
			t.Fatal("ожидалась ошибка")
		}

		if event != nil {
			t.Errorf("ожидалось nil-событие, получено %#v", event)
		}

		if !errors.Is(err, interfaces.ErrNoEvent) {
			t.Fatalf("ожидалась ErrNoEvent, получено: %v", err)
		}
	})

	t.Run("zero ID", func(t *testing.T) {
		event := newEvent("event with zero ID", false)

		deleted, err := storage.DeleteEvent(ctx, event)
		if err == nil {
			t.Fatal("ожидалась ошибка")
		}

		if deleted != event {
			t.Error("ожидался исходный объект события")
		}

		if !errors.Is(err, interfaces.ErrNoEvent) {
			t.Fatalf("ожидалась ErrNoEvent, получено: %v", err)
		}
	})
}

func TestStorage_ListEvents(t *testing.T) {
	storage := newTestStorage(t)
	ctx := context.Background()

	source := []*models.Event{
		newEvent("event-1", false),
		newEvent("event-2", true),
		newEvent("event-3", false),
	}

	expectedByID := make(map[int]*models.Event, len(source))

	for _, event := range source {
		created, err := storage.CreateEvent(ctx, event)
		if err != nil {
			t.Fatalf("не удалось создать событие %q: %v", event.Title, err)
		}

		expectedByID[created.ID] = created
	}

	events, err := storage.ListEvents(ctx)
	if err != nil {
		t.Fatalf("ListEvents вернул ошибку: %v", err)
	}

	if len(events) != len(source) {
		t.Fatalf(
			"ожидалось %d событий, получено %d",
			len(source),
			len(events),
		)
	}

	for i := range events {
		expected, ok := expectedByID[events[i].ID]
		if !ok {
			t.Errorf("получено неожиданное событие с ID=%d", events[i].ID)
			continue
		}

		assertEventsEqual(t, expected, &events[i])
	}
}

func TestStorage_ListNotSheduledEvents(t *testing.T) {
	storage := newTestStorage(t)
	ctx := context.Background()

	source := []*models.Event{
		newEvent("not-sheduled-1", false),
		newEvent("sheduled", true),
		newEvent("not-sheduled-2", false),
	}

	expectedTitles := map[string]bool{
		"not-sheduled-1": true,
		"not-sheduled-2": true,
	}

	for _, event := range source {
		if _, err := storage.CreateEvent(ctx, event); err != nil {
			t.Fatalf("не удалось создать событие %q: %v", event.Title, err)
		}
	}

	events, err := storage.ListNotSheduledEvents(ctx)
	if err != nil {
		t.Fatalf("ListNotSheduledEvents вернул ошибку: %v", err)
	}

	if len(events) != len(expectedTitles) {
		t.Fatalf(
			"ожидалось %d незапланированных событий, получено %d",
			len(expectedTitles),
			len(events),
		)
	}

	for _, event := range events {
		if event.Sheduled {
			t.Errorf(
				"ListNotSheduledEvents вернул запланированное событие ID=%d",
				event.ID,
			)
		}

		if !expectedTitles[event.Title] {
			t.Errorf("получено неожиданное событие %q", event.Title)
		}
	}
}

func TestStorage_ListEventsEmpty(t *testing.T) {
	storage := newTestStorage(t)

	events, err := storage.ListEvents(context.Background())
	if err != nil {
		t.Fatalf("ListEvents вернул ошибку: %v", err)
	}

	if len(events) != 0 {
		t.Fatalf("ожидался пустой список, получено %d событий", len(events))
	}
}

func TestStorage_CloseWithoutConnect(t *testing.T) {
	storage := New("pgx", defaultTestDSN)

	if err := storage.Close(); err != nil {
		t.Fatalf("Close без Connect вернул ошибку: %v", err)
	}
}
