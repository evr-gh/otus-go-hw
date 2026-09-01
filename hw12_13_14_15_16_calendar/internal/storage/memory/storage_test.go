package memorystorage

import (
	"context"
	"testing"
	"time"

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

func TestStorage(t *testing.T) {
	cntx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage := New()
	err := storage.Connect(cntx)
	require.NoError(t, err)

	defer func() {
		err := storage.Close(cntx)
		require.NoError(t, err)
	}()

	// CREATE

	d1, err := time.ParseDuration("5m7s")
	require.NoError(t, err)
	event1.Duration = d1

	nlt1, err := time.ParseDuration("10m")
	require.NoError(t, err)
	event1.NotifyLeadTime = nlt1

	firstEvent, err := storage.CreateEvent(cntx, &event1)
	require.NoError(t, err)
	require.Equal(t, 1, firstEvent.ID)

	d2, err := time.ParseDuration("10s")
	require.NoError(t, err)
	event2.Duration = d2

	nlt2, err := time.ParseDuration("2h")
	require.NoError(t, err)
	event2.NotifyLeadTime = nlt2

	secondEvent, err := storage.CreateEvent(cntx, &event2)
	require.NoError(t, err)
	require.Equal(t, 2, secondEvent.ID)

	// DELETE
	dEvent2, err := storage.DeleteEvent(cntx, secondEvent)
	require.NoError(t, err)
	require.Equal(t, secondEvent, dEvent2)

	dEvent3, err := storage.DeleteEvent(cntx, secondEvent)
	require.Error(t, err)
	require.Equal(t, "не удалено событие с ID=2: нет в БД", err.Error())
	require.Equal(t, secondEvent, dEvent3)

	// CREATE
	d3, err := time.ParseDuration("1m")
	require.NoError(t, err)
	event3.Duration = d3

	nlt3, err := time.ParseDuration("2m1s")
	require.NoError(t, err)
	event3.NotifyLeadTime = nlt3

	thirdEvent, err := storage.CreateEvent(cntx, &event3)
	require.NoError(t, err)
	require.Equal(t, 3, thirdEvent.ID)

	// UPDATE
	description := "Обновленное описание третьего события"
	thirdEvent.Description = description
	thirdEvent.Sheduled = true
	uEvent1, err := storage.UpdateEvent(cntx, thirdEvent)
	require.NoError(t, err)
	require.Equal(t, thirdEvent, uEvent1)
	require.Equal(t, description, uEvent1.Description)

	uEvent2, err := storage.UpdateEvent(cntx, secondEvent)
	require.Error(t, err)
	require.Equal(t, "не обновлено событие с ID=2: нет в БД", err.Error())
	require.Equal(t, secondEvent, uEvent2)

	// READ (all)
	events, err := storage.ListEvents(cntx)
	require.NoError(t, err)

	require.Equal(t, 2, len(events))

	tevents := make([]data.Event, 0, 2)

	for _, e := range events {
		if e.ID == 1 {
			tevents = append(tevents, e)
		}
	}
	for _, e := range events {
		if e.ID == 3 {
			tevents = append(tevents, e)
		}
	}

	require.Equal(t, []data.Event{*firstEvent, *thirdEvent}, tevents)

	// READ
	rEvent, err := storage.ReadEvent(cntx, 3)
	require.NoError(t, err)
	require.Equal(t, thirdEvent, rEvent)

	rEvent, err = storage.ReadEvent(cntx, 4)
	require.Error(t, err)
	require.Equal(t, "не получено событие с ID=4: нет в БД", err.Error())
	require.Nil(t, rEvent)

	// READ
	nsEvents, err := storage.ListNotSheduledEvents(cntx)
	require.NoError(t, err)

	require.Equal(t, 1, len(nsEvents))
	require.Equal(t, []data.Event{*firstEvent}, nsEvents)
}
