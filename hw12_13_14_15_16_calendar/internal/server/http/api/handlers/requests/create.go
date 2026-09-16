package api

import (
	"encoding/json"
	"net/http"

	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
	models "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/models"
	responses "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/api/api_response"
	commonHandlers "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/api/handlers/common"
)

type EventsCreateHandler struct {
	Logger interfaces.Logger
	App    interfaces.Application
}

func (h EventsCreateHandler) ServeHTTP(rw http.ResponseWriter, rr *http.Request) {
	apiMethod := "api.events.create"
	if rr == nil {
		commonHandlers.InvalidRequestBodyHandler{APIMethod: apiMethod}.ServeHTTP(rw, rr)
	}
	if rr.Method != commonHandlers.POST {
		commonHandlers.InvalidHTTPMethod{APIMethod: apiMethod}.ServeHTTP(rw, rr)
		return
	}
	var event models.Event
	err := json.NewDecoder(rr.Body).Decode(&event)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	createdEvent, err := h.App.CreateEvent(rr.Context(), &event)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	apiResponse := responses.NewAPIResponse(apiMethod)
	apiResponse.Data = responses.DataItem{Item: createdEvent}
	apiResponseJSON, err := json.Marshal(apiResponse)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(apiResponseJSON)
}
