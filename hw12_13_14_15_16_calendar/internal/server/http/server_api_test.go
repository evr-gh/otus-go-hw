package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	app "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	logger "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	models "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/models"
	middleware "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/middleware"
	storage "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

type APIResponseTest struct {
	APIMethod string
	Error     string
	Data      struct{ Item models.Event }
}

type APIMultyResponseTest struct {
	APIMethod string
	Error     string
	Data      struct{ Items []models.Event }
}

//nolint:funlen
func TestServerAPI(t *testing.T) {
	var port uint16 = 5011
	host := "localhost"
	var response *http.Response
	var err error
	var apiResponse APIResponseTest
	ctx := context.Background()

	outputInto := &bytes.Buffer{}
	logg := logger.New(logger.INFO, outputInto)
	middleware.Init(logg)

	calendarApp := app.New(logg, storage.New())

	httpServer := NewHTTPServer(
		calendarApp,
		host,
		port,
		10*time.Second,
		12*time.Second,
		14*time.Second,
		65000,
		logg,
	)
	wg := sync.WaitGroup{}
	wg.Go(func() {
		err := httpServer.Start(ctx)
		require.NoErrorf(t, err, "Ошибка при запуске сервера")
	})
	time.Sleep(2 * time.Second)
	//
	client := &http.Client{}

	requestOfCreate := fmt.Sprintf("http://%s:%d/api/events/create", host, port)
	payloadOfCreateRaw := `{"title": "title 1", "startat": "2023-08-05T21:54:42+02:00"}`

	payloadOfCreate := strings.NewReader(payloadOfCreateRaw)
	request, err := http.NewRequestWithContext(ctx, "POST", requestOfCreate, payloadOfCreate)
	require.NoErrorf(t, err, "Не удалось создать запрос на создание события 1")
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	require.NoErrorf(t, err, "Не удалось выполнить запрос на создание события 1")
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	require.NoErrorf(t, err, "Не удалось декодировать ответ на запрос на создание события 1")
	require.Equal(t, "", apiResponse.Error, "Получена ошибка в ответ на запрос на создание события 1")
	require.Equal(t, 1, apiResponse.Data.Item.ID, "Неверный ID в ответе на запрос на создание события 1")
	require.Equal(t, "title 1", apiResponse.Data.Item.Title, "Неверный Title в ответе на запрос на создание события 1")
	response.Body.Close()

	payloadOfCreateRaw = `{"title": "title 2", "startat": "2023-08-05T21:54:45+02:00"}`
	payloadOfCreate = strings.NewReader(payloadOfCreateRaw)
	request, err = http.NewRequestWithContext(ctx, "POST", requestOfCreate, payloadOfCreate)
	require.NoErrorf(t, err, "Не удалось создать запрос на создание события 2")
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	require.NoErrorf(t, err, "Не удалось выполнить запрос на создание события 2")
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	require.NoErrorf(t, err, "Не удалось декодировать ответ на запрос на создание события 2")
	require.Equal(t, "", apiResponse.Error, "Получена ошибка в ответ на запрос на создание события 2")
	require.Equal(t, 2, apiResponse.Data.Item.ID, "Неверный ID в ответе на запрос на создание события 2")
	response.Body.Close()

	payloadOfCreateRaw = `{"title": "title 3", "startat": "2023-08-05T21:54:55+02:00"}`
	payloadOfCreate = strings.NewReader(payloadOfCreateRaw)
	request, err = http.NewRequestWithContext(ctx, "POST", requestOfCreate, payloadOfCreate)
	require.NoErrorf(t, err, "Не удалось создать запрос на создание события 3")
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	require.NoErrorf(t, err, "Не удалось выполнить запрос на создание события 3")
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	require.NoErrorf(t, err, "Не удалось декодировать ответ на запрос на создание события 3")
	require.Equal(t, "", apiResponse.Error, "Получена ошибка в ответ на запрос на создание события 3")
	require.Equal(t, 3, apiResponse.Data.Item.ID, "Неверный ID в ответе на запрос на создание события 3")
	response.Body.Close()

	requestOfDelete := fmt.Sprintf("http://%s:%d/api/events/2/delete", host, port)
	payloadOfDelete := strings.NewReader("")
	request, err = http.NewRequestWithContext(ctx, "DELETE", requestOfDelete, payloadOfDelete)
	require.NoErrorf(t, err, "Не удалось создать запрос на удаление события 2")
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	require.NoErrorf(t, err, "Не удалось выполнить запрос на удаление события 2")
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	require.NoErrorf(t, err, "Не удалось декодировать ответ на запрос на удаление события 2")
	require.Equal(t, "", apiResponse.Error, "Получена ошибка в ответ на запрос на удаление события 2")
	response.Body.Close()

	requestOfUpdate := fmt.Sprintf("http://%s:%d/api/events/3/update", host, port)
	payloadOfUpdateRaw := `{"title": "title 3 updated", "startat": "2023-08-05T21:54:55+02:00", "sheduled":true}`
	payloadOfUpdate := strings.NewReader(payloadOfUpdateRaw)
	request, err = http.NewRequestWithContext(ctx, "PATCH", requestOfUpdate, payloadOfUpdate)
	require.NoErrorf(t, err, "Не удалось создать запрос на обновление события 3")
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	require.NoErrorf(t, err, "Не удалось выполнить запрос на обновление события 3")
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	require.NoErrorf(t, err, "Не удалось декодировать ответ на запрос на обновление события 3")
	require.Equal(t, "", apiResponse.Error, "Получена ошибка в ответ на запрос на обновление события 3")
	require.Equal(t, 3, apiResponse.Data.Item.ID, "Неверный ID в ответе на запрос на обновление события 3")
	require.Equal(t, "title 3 updated", apiResponse.Data.Item.Title, "Неверный Title в ответе"+
		" на запрос на обновление события 3")
	require.Equal(t, true, apiResponse.Data.Item.Sheduled, "Неверный Sheduled в ответе на запрос на обновление события 3")
	response.Body.Close()

	requestOfGet := fmt.Sprintf("http://%s:%d/api/events/3", host, port)
	payloadOfGet := strings.NewReader("")
	request, err = http.NewRequestWithContext(ctx, "GET", requestOfGet, payloadOfGet)
	require.NoErrorf(t, err, "Не удалось создать запрос на получение события 3")
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	require.NoErrorf(t, err, "Не удалось выполнить запрос на получение события 3")
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	require.NoErrorf(t, err, "Не удалось декодировать ответ на запрос на получение события 3")
	require.Equal(t, "", apiResponse.Error, "Получена ошибка в ответ на запрос на получение события 3")
	require.Equal(t, 3, apiResponse.Data.Item.ID, "Неверный ID в ответе на запрос на получение события 3")
	require.Equal(t, "title 3 updated", apiResponse.Data.Item.Title, "Неверный Title в ответе"+
		" на запрос на получение события 3")
	require.Equal(t, true, apiResponse.Data.Item.Sheduled, "Неверный Sheduled в ответе на запрос на получение события 3")
	response.Body.Close()

	var apiMultyResponse1 APIMultyResponseTest
	payloadOfList := strings.NewReader(``)
	requestOfList := fmt.Sprintf("http://%s:%d/api/events/", host, port)
	request, err = http.NewRequestWithContext(ctx, "GET", requestOfList, payloadOfList)
	require.NoErrorf(t, err, "Не удалось создать запрос на получение списка событий")
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	require.NoErrorf(t, err, "Не удалось выполнить запрос на получение списка событий")
	err = json.NewDecoder(response.Body).Decode(&apiMultyResponse1)
	require.NoErrorf(t, err, "Не удалось декодировать ответ на запрос на получение списка событий")
	response.Body.Close()
	require.Equal(t, 2, len(apiMultyResponse1.Data.Items), "Неверное количестов элементов"+
		" в ответе на запрос на получение списка событий")

	var apiMultyResponse2 APIMultyResponseTest
	payloadOfListNotSheduled := strings.NewReader(``)
	requestOfListNotSheduled := fmt.Sprintf("http://%s:%d/api/events/notsheduled", host, port)
	request, err = http.NewRequestWithContext(ctx, "GET", requestOfListNotSheduled, payloadOfListNotSheduled)
	require.NoErrorf(t, err, "Не удалось создать запрос на получение списка незапланированных событий")
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	require.NoErrorf(t, err, "Не удалось выполнить запрос на получение списка незапланированных событий")
	err = json.NewDecoder(response.Body).Decode(&apiMultyResponse2)
	require.NoErrorf(t, err, "Не удалось декодировать ответ на запрос на получение списка незапланированных событий")
	response.Body.Close()
	require.Equal(t, 1, len(apiMultyResponse2.Data.Items), "Неверное количестов элементов в ответе на запрос на"+
		" получение списка незапланированных событий")
	require.Equal(t, 1, apiMultyResponse2.Data.Items[0].ID, "Неверный ID в ответе на запрос на получение списка"+
		" незапланированных событий")
	require.Equal(t, "title 1", apiMultyResponse2.Data.Items[0].Title, "Неверный Title в ответе на запрос"+
		" на получение списка незапланированных событий")

	httpServer.Stop(context.Background())
	wg.Wait()
}

func TestServerAPIVersion(t *testing.T) {
	var port uint16 = 5012
	host := "localhost"
	ctx := context.Background()

	outputInto := &bytes.Buffer{}
	logg := logger.New(logger.INFO, outputInto)
	middleware.Init(logg)
	calendarApp := app.New(logg, storage.New())

	httpServer := NewHTTPServer(
		calendarApp,
		host,
		port,
		10*time.Second,
		12*time.Second,
		14*time.Second,
		65000,
		logg,
	)
	wg := sync.WaitGroup{}
	wg.Go(func() {
		err := httpServer.Start(ctx)
		require.NoErrorf(t, err, "Ошибка при запуске сервера")
	})

	time.Sleep(2 * time.Second)

	client := &http.Client{}
	requestOfVersion := fmt.Sprintf("http://%s:%d/api/version", host, port)
	request, err := http.NewRequestWithContext(context.Background(), "GET", requestOfVersion, strings.NewReader(``))
	require.NoErrorf(t, err, "Не удалось создать запрос версии")
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	require.NoErrorf(t, err, "Не удалось выполнить запрос версии")
	body, err := io.ReadAll(response.Body)
	require.NoErrorf(t, err, "Не удалось считать ответ на запрос версии")
	response.Body.Close()
	ethalon := `{"method":"api.version","error":"","data":{"Version":"1.0.0"}}`
	require.Equal(t, ethalon, string(body), "Неверный ответ на запрос версии")

	httpServer.Stop(context.Background())
	wg.Wait()
}
