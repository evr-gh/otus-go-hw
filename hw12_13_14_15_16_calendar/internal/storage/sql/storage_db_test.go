package sqlstorage

import (
	"context"
	"os"
	"testing"
	"time"

	models "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/data"
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
		Duration:       5_400_000_000_000, // 90 минут в наносекундах.
		Owner:          "integration-test-user",
		NotifyLeadTime: 900_000_000_000, // 15 минут в наносекундах.
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
