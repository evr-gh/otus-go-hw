package sqlstorage

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	data "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/data"
	"github.com/stretchr/testify/require"
)

var event1 = data.Event{
	ID:          0,
	Title:       "Title1",
	Time:        time.Now(),
	Description: "Описание первого события",
	Owner:       "tester",
}

var event2 = data.Event{
	ID:          0,
	Title:       "Title2",
	Time:        time.Now(),
	Description: "Описание второго события",
	Owner:       "tester",
}

var event3 = data.Event{
	ID:          0,
	Title:       "Title3",
	Time:        time.Now(),
	Description: "Описание третьего события",
	Owner:       "tester2",
}

//nolint:funlen
func TestStorage(t *testing.T) {
	cntx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const dsn = "test_db"
	stdDB, mock, err := sqlmock.NewWithDSN(dsn)
	require.NoError(t, err)
	defer stdDB.Close()

	storage := New("sqlmock", dsn)
	err = storage.Connect(cntx)
	require.NoError(t, err)

	d1, err := time.ParseDuration("5m7s")
	require.NoError(t, err)
	event1.Duration = d1

	nlt1, err := time.ParseDuration("10m")
	require.NoError(t, err)
	event1.NotifyLeadTime = nlt1

	mock.
		ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO events
			("title", "description", "time", "duration", "owner", "notifyleadtime", "sheduled")
		values($1, $2, $3, $4, $5, $6, $7) RETURNING "id";
	`)).
		WithArgs(
			event1.Title,
			event1.Description,
			event1.Time,
			event1.Duration,
			event1.Owner,
			event1.NotifyLeadTime,
			event1.Sheduled,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(int64(1)),
		)

	d2, err := time.ParseDuration("10s")
	require.NoError(t, err)
	event2.Duration = d2

	nlt2, err := time.ParseDuration("2h")
	require.NoError(t, err)
	event2.NotifyLeadTime = nlt2

	mock.
		ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO events
			("title", "description", "time", "duration", "owner", "notifyleadtime", "sheduled")
		values($1, $2, $3, $4, $5, $6, $7) RETURNING "id";
	`)).
		WithArgs(
			event2.Title,
			event2.Description,
			event2.Time,
			event2.Duration,
			event2.Owner,
			event2.NotifyLeadTime,
			event2.Sheduled,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(int64(2)),
		)

	// CREATE

	firstEvent, err := storage.CreateEvent(cntx, &event1)
	require.NoError(t, err)
	require.Equal(t, 1, firstEvent.ID)

	secondEvent, err := storage.CreateEvent(cntx, &event2)
	require.NoError(t, err)
	require.Equal(t, 2, secondEvent.ID)

	nilEvent, err := storage.CreateEvent(cntx, nil)
	require.Error(t, err)
	require.Equal(t, "не создано событие: не передана информация по событию", err.Error())
	require.Nil(t, nilEvent)

	// DELETE

	mock.
		ExpectExec(regexp.QuoteMeta(`DELETE FROM events WHERE id=$1;`)).
		WithArgs(secondEvent.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	dEvent2, err := storage.DeleteEvent(cntx, secondEvent)
	require.NoError(t, err)
	require.Equal(t, secondEvent, dEvent2)

	secondEvent.ID = 2

	mock.
		ExpectExec(regexp.QuoteMeta(`DELETE FROM events WHERE id=$1;`)).
		WithArgs(secondEvent.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	dEvent3, err := storage.DeleteEvent(cntx, secondEvent)
	require.NoError(t, err)
	require.Equal(t, secondEvent, dEvent3)

	nilEvent, err = storage.DeleteEvent(cntx, nil)
	require.Error(t, err)
	require.Equal(t, "не удалено событие: не передана информация по событию", err.Error())
	require.Nil(t, nilEvent)

	d3, err := time.ParseDuration("1m")
	require.NoError(t, err)
	event3.Duration = d3

	nlt3, err := time.ParseDuration("2m1s")
	require.NoError(t, err)
	event3.NotifyLeadTime = nlt3

	mock.
		ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO events
			("title", "description", "time", "duration", "owner", "notifyleadtime", "sheduled")
		values($1, $2, $3, $4, $5, $6, $7) RETURNING "id";
	`)).
		WithArgs(
			event3.Title,
			event3.Description,
			event3.Time,
			event3.Duration,
			event3.Owner,
			event3.NotifyLeadTime,
			event3.Sheduled,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(int64(3)),
		)

	thirdEvent, err := storage.CreateEvent(cntx, &event3)
	require.NoError(t, err)
	require.Equal(t, 3, thirdEvent.ID)

	description := "Обновленное описание третьего события"
	thirdEvent.Description = description
	thirdEvent.Sheduled = true

	mock.
		ExpectExec(regexp.QuoteMeta(` UPDATE events SET "title"=$1, "description"=$2, "time"=$3, "duration"=$4, 
	"owner"=$5, "notifyleadtime"=$6, "sheduled"=$7 WHERE id=$8;`)).
		WithArgs(
			event3.Title,
			event3.Description,
			event3.Time,
			event3.Duration,
			event3.Owner,
			event3.NotifyLeadTime,
			event3.Sheduled,
			event3.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	uEvent1, err := storage.UpdateEvent(cntx, thirdEvent)
	require.NoError(t, err)
	require.Equal(t, thirdEvent, uEvent1)
	require.Equal(t, description, uEvent1.Description)

	secondEvent.ID = 2

	mock.
		ExpectExec(regexp.QuoteMeta(` UPDATE events SET "title"=$1, "description"=$2, "time"=$3, "duration"=$4, 
	"owner"=$5, "notifyleadtime"=$6, "sheduled"=$7 WHERE id=$8;`)).
		WithArgs(
			secondEvent.Title,
			secondEvent.Description,
			secondEvent.Time,
			secondEvent.Duration,
			secondEvent.Owner,
			secondEvent.NotifyLeadTime,
			secondEvent.Sheduled,
			secondEvent.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	uEvent2, err := storage.UpdateEvent(cntx, secondEvent)
	require.Error(t, err)
	require.Equal(t, "не обновлено событие с ID=2: нет в БД", err.Error())
	require.Equal(t, secondEvent, uEvent2)

	nilEvent, err = storage.UpdateEvent(cntx, nil)
	require.Error(t, err)
	require.Equal(t, "не обновлено событие: не передана информация по событию", err.Error())
	require.Nil(t, nilEvent)

	fieldNames := []string{"id", "title", "description", "time", "duration", "owner", "notifyleadtime", "sheduled"}

	sqlReadEvent := `SELECT "id", "title", "description", "time", "duration", "owner", "notifyleadtime", "sheduled"
	FROM events WHERE "id"=$1;`

	mock.
		ExpectQuery(regexp.QuoteMeta(sqlReadEvent)).
		WithArgs(
			3,
		).
		WillReturnRows(
			sqlmock.NewRows(fieldNames).
				AddRow(thirdEvent.ID, thirdEvent.Title, thirdEvent.Description, thirdEvent.Time,
					thirdEvent.Duration, thirdEvent.Owner, thirdEvent.NotifyLeadTime, thirdEvent.Sheduled),
		)

	rEvent, err := storage.ReadEvent(cntx, 3)
	require.NoError(t, err)
	require.Equal(t, thirdEvent, rEvent)

	mock.
		ExpectQuery(regexp.QuoteMeta(sqlReadEvent)).
		WithArgs(
			4,
		).
		WillReturnRows(
			sqlmock.NewRows(fieldNames),
		)
	rEvent, err = storage.ReadEvent(cntx, 4)
	require.Error(t, err)
	require.Equal(t, "не получено событие с ID=4: нет в БД", err.Error())
	require.Nil(t, rEvent)

	sqlListEvents := `SELECT "id", "title", "description", "startat", "durationseconds", "owner", 
	"notifyearlyseconds", "sheduled" FROM events;`

	mock.
		ExpectQuery(regexp.QuoteMeta(sqlListEvents)).
		WillReturnRows(
			sqlmock.NewRows(fieldNames).
				AddRow(firstEvent.ID, firstEvent.Title, firstEvent.Description, firstEvent.Time,
					firstEvent.Duration, firstEvent.Owner, firstEvent.NotifyLeadTime, firstEvent.Sheduled).
				AddRow(thirdEvent.ID, thirdEvent.Title, thirdEvent.Description, thirdEvent.Time,
					thirdEvent.Duration, thirdEvent.Owner, thirdEvent.NotifyLeadTime, thirdEvent.Sheduled),
		)

	events, err := storage.ListEvents(cntx)
	require.NoError(t, err)
	require.Equal(t, 2, len(events))
	require.Equal(t, []data.Event{*firstEvent, *thirdEvent}, events)

	mock.
		ExpectQuery(regexp.QuoteMeta(sqlListEvents)).
		WillReturnRows(
			sqlmock.NewRows(fieldNames).
				AddRow("ID", firstEvent.Title, firstEvent.Description, firstEvent.Time,
					firstEvent.Duration, firstEvent.Owner, firstEvent.NotifyLeadTime, firstEvent.Sheduled),
		)
	events, err = storage.ListEvents(cntx)
	require.Error(t, err)
	require.Equal(t, "список событий не получен: sql:"+
		" Scan error on column index 0, name \"id\":"+
		" converting driver.Value type string (\"ID\") to a int:"+
		" invalid syntax", err.Error())
	require.Nil(t, events)

	mock.
		ExpectQuery(regexp.QuoteMeta(sqlListEvents)).
		WillReturnError(fmt.Errorf("db connection error"))

	events, err = storage.ListEvents(cntx)
	require.Error(t, err)
	require.Equal(t, "список событий не получен: db connection error", err.Error())
	require.Nil(t, events)

	sqlListNotSheduledEvents := `SELECT "id", "title", "description", "startat", "durationseconds", "owner", 
	"notifyearlyseconds", "sheduled" FROM events WHERE "sheduled" IS NOT TRUE;`

	mock.
		ExpectQuery(regexp.QuoteMeta(sqlListNotSheduledEvents)).
		WillReturnRows(
			sqlmock.NewRows(fieldNames).
				AddRow(firstEvent.ID, firstEvent.Title, firstEvent.Description, firstEvent.Time,
					firstEvent.Duration, firstEvent.Owner, firstEvent.NotifyLeadTime, firstEvent.Sheduled),
		)

	events, err = storage.ListNotSheduledEvents(cntx)
	require.NoError(t, err)
	require.Equal(t, 1, len(events))
	require.Equal(t, []data.Event{*firstEvent}, events)

	mock.
		ExpectQuery(regexp.QuoteMeta(sqlListNotSheduledEvents)).
		WillReturnRows(
			sqlmock.NewRows(fieldNames).
				AddRow("ID", firstEvent.Title, firstEvent.Description, firstEvent.Time,
					firstEvent.Duration, firstEvent.Owner, firstEvent.NotifyLeadTime, firstEvent.Sheduled),
		)
	events, err = storage.ListNotSheduledEvents(cntx)
	require.Error(t, err)
	require.Equal(t, "список незапланированных событий не получен: sql:"+
		" Scan error on column index 0, name \"id\":"+
		" converting driver.Value type string (\"ID\") to a int:"+
		" invalid syntax", err.Error())
	require.Nil(t, events)

	mock.
		ExpectQuery(regexp.QuoteMeta(sqlListNotSheduledEvents)).
		WillReturnError(fmt.Errorf("db connection error"))

	events, err = storage.ListNotSheduledEvents(cntx)
	require.Error(t, err)
	require.Equal(t, "список незапланированных событий не получен: db connection error", err.Error())
	require.Nil(t, events)

	mock.ExpectClose()

	err = storage.Close()
	require.NoError(t, err)
}
