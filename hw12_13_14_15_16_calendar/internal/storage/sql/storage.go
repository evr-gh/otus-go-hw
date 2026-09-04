package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	models "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/data"
	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
	// mysql driver.
	_ "github.com/go-sql-driver/mysql"
	// Postgersql driver.
	_ "github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	// sqlite driver.
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	SQLDbType string
	DSN       string
	db        *sqlx.DB
}

func New(sqlDbType string, dsn string) *Storage {
	return &Storage{sqlDbType, dsn, nil}
}

func (s *Storage) Connect(cntx context.Context) error {
	db, err := sqlx.ConnectContext(cntx, s.SQLDbType, s.DSN)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к БД: %w", err)
	}
	s.db = db
	return nil
}

func (s *Storage) Close() error {
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("возникла ошибка при отключении от БД: %w", err)
	}
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error) {
	if event == nil {
		return event, fmt.Errorf("не создано событие: %w", interfaces.ErrNoEvent)
	}
	sqlStatement := `INSERT INTO events 
			("title", "description", "time", "duration", "owner", "notifyleadtime", "sheduled") 
			values($1, $2, $3, $4, $5, $6, $7) RETURNING "id";`
	err := s.db.QueryRowxContext(ctx, sqlStatement,
		event.Title, event.Description, event.Time, event.Duration, event.Owner, event.NotifyLeadTime, event.Sheduled,
	).Scan(&event.ID)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать событие %q (%v): %w", event.Title, event.Time, err)
	}
	return event, nil
}

func (s *Storage) ReadEvent(ctx context.Context, eventID int) (*models.Event, error) {
	var e models.Event
	sqlStatement := `SELECT "id", "title", "description", "time", "duration", "owner", "notifyleadtime", "sheduled"
	FROM events WHERE "id"=$1;`
	err := s.db.QueryRowxContext(ctx, sqlStatement, eventID).Scan(&e.ID,
		&e.Title, &e.Description, &e.Time, &e.Duration,
		&e.Owner, &e.NotifyLeadTime, &e.Sheduled)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("не получено событие с ID=%v: %w", eventID, interfaces.ErrNoData)
		}
		return nil, fmt.Errorf("не получено событие с ID=%v: %w", eventID, err)
	}
	return &e, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event *models.Event) (*models.Event, error) {
	if event == nil {
		return event, fmt.Errorf("не обновлено событие: %w", interfaces.ErrNoEvent)
	}
	sqlStatement := `UPDATE events SET "title"=$1, "description"=$2, "time"=$3, "duration"=$4, 
	"owner"=$5, "notifyleadtime"=$6, "sheduled"=$7 WHERE id=$8;`
	res, err := s.db.ExecContext(ctx, sqlStatement, event.Title, event.Description,
		event.Time, event.Duration, event.Owner, event.NotifyLeadTime, event.Sheduled,
		event.ID)
	if err != nil {
		return event, fmt.Errorf("не обновлено событие: %w", err)
	}
	if irows, err := res.RowsAffected(); (irows == int64(0)) || (err != nil) {
		if err != nil {
			return event, fmt.Errorf("не обновлено событие с ID=%v: %w", event.ID, err)
		}
		return event, fmt.Errorf("не обновлено событие с ID=%v: %w", event.ID, interfaces.ErrNoData)
	}
	return event, nil
}

func (s *Storage) DeleteEvent(ctx context.Context, event *models.Event) (*models.Event, error) {
	if event == nil {
		return event, fmt.Errorf("не удалено событие: %w", interfaces.ErrNoEvent)
	}
	if event.ID == 0 {
		return event, fmt.Errorf("не удалено событие: %w", interfaces.ErrNoEvent)
	}
	sqlStatement := `DELETE FROM events WHERE id=$1;`
	_, err := s.db.ExecContext(ctx, sqlStatement, event.ID)
	if err != nil {
		return nil, fmt.Errorf("не удалено событие c ID=%v: %w", event.ID, err)
	}

	event.ID = 0
	return event, nil
}

func (s *Storage) ListEvents(ctx context.Context) ([]models.Event, error) {
	var e models.Event
	var events []models.Event
	sqlStatement := `SELECT "id", "title", "description", "startat", "durationseconds", "owner", 
	"notifyearlyseconds", "sheduled" FROM events;`
	rows, err := s.db.QueryContext(ctx, sqlStatement)
	fmt.Println(rows)
	if err != nil {
		return nil, fmt.Errorf("список событий не получен: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&e.ID, &e.Title, &e.Description, &e.Time, &e.Duration, &e.Owner, &e.NotifyLeadTime, &e.Sheduled)
		if err != nil {
			return nil, fmt.Errorf("список событий не получен: %w", err)
		}
		events = append(events, e)
	}
	return events, nil
}

func (s *Storage) ListNotSheduledEvents(ctx context.Context) ([]models.Event, error) {
	var e models.Event
	var events []models.Event
	sqlStatement := `SELECT "id", "title", "description", "startat", "durationseconds", "owner", 
	"notifyearlyseconds", "sheduled" FROM events WHERE "sheduled" IS NOT TRUE;`
	rows, err := s.db.QueryContext(ctx, sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("список незапланированных событий не получен: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&e.ID, &e.Title, &e.Description, &e.Time, &e.Duration, &e.Owner, &e.NotifyLeadTime, &e.Sheduled)
		if err != nil {
			return nil, fmt.Errorf("список незапланированных событий не получен: %w", err)
		}
		events = append(events, e)
	}
	return events, nil
}
