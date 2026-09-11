package api

import (
	"encoding/json"
	"net/http"

	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
	models "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/models"
	responses "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/api/api_response"
	commonHandlers "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/api/handlers/common"
)

type EventsListHandler struct {
	Logger    interfaces.Logger
	App       interfaces.Application
	APIMethod string
}

func (h EventsListHandler) ServeHTTP(rw http.ResponseWriter, rr *http.Request) {
	if rr == nil {
		commonHandlers.InvalidRequestBodyHandler{APIMethod: h.APIMethod}.ServeHTTP(rw, rr)
	}
	if rr.Method != commonHandlers.GET {
		commonHandlers.InvalidHTTPMethod{APIMethod: h.APIMethod}.ServeHTTP(rw, rr)
		return
	}
	var events []models.Event
	var err error
	if h.APIMethod == "api.events.listnotsheduled" {
		events, err = h.App.ListNotSheduledEvents(rr.Context())
	} else {
		events, err = h.App.ListEvents(rr.Context())
	}
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: h.APIMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	apiResponse := responses.NewAPIResponse(h.APIMethod)
	apiResponse.Data = responses.DataItems{Items: events}
	apiResponseJSON, err := json.Marshal(apiResponse)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: h.APIMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(apiResponseJSON)
}
