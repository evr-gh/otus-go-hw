package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
	models "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/models"
	responses "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/api/api_response"
	commonHandlers "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/api/handlers/common"
)

type EventsDeleteHandler struct {
	Logger interfaces.Logger
	App    interfaces.Application
}

func (h EventsDeleteHandler) ServeHTTP(rw http.ResponseWriter, rr *http.Request) {
	apiMethod := "api.events.delete"
	if rr == nil {
		commonHandlers.InvalidRequestBodyHandler{APIMethod: apiMethod}.ServeHTTP(rw, rr)
	}
	if rr.Method != "DELETE" {
		commonHandlers.InvalidHTTPMethod{APIMethod: apiMethod}.ServeHTTP(rw, rr)
		return
	}
	idString := rr.Form.Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	event := models.Event{}
	event.ID = id
	deletedEvent, err := h.App.DeleteEvent(rr.Context(), &event)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	apiResponse := responses.NewAPIResponse(apiMethod)
	apiResponse.Data = responses.DataItem{Item: deletedEvent}
	apiResponseJSON, err := json.Marshal(apiResponse)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(apiResponseJSON)
}
