package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
	responses "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/api/api_response"
	commonHandlers "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/api/handlers/common"
)

type EventsGetHandler struct {
	Logger interfaces.Logger
	App    interfaces.Application
}

func (h EventsGetHandler) ServeHTTP(rw http.ResponseWriter, rr *http.Request) {
	apiMethod := "api.events.get"
	if rr == nil {
		commonHandlers.InvalidRequestBodyHandler{APIMethod: apiMethod}.ServeHTTP(rw, rr)
	}
	if rr.Method != commonHandlers.GET {
		commonHandlers.InvalidHTTPMethod{APIMethod: apiMethod}.ServeHTTP(rw, rr)
		return
	}
	idString := rr.Form.Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	event, err := h.App.ReadEvent(rr.Context(), id)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	apiResponse := responses.NewAPIResponse(apiMethod)
	apiResponse.Data = responses.DataItem{Item: event}
	apiResponseJSON, err := json.Marshal(apiResponse)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(apiResponseJSON)
}
