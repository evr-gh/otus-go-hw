package api

import (
	"encoding/json"
	"net/http"

	responses "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/api/api_response"
)

var (
	GET  = "GET"
	POST = "POST"
)

type InvalidHTTPMethod struct {
	APIMethod string
}

func (h InvalidHTTPMethod) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	apiResponse := responses.InvalidHTTPMethod(h.APIMethod)
	apiResponseJSON, _ := json.Marshal(apiResponse)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusMethodNotAllowed)
	rw.Write(apiResponseJSON)
}

type InvalidRequestBodyHandler struct {
	APIMethod string
}

func (h InvalidRequestBodyHandler) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	apiResponse := responses.InvalidRequestBody(h.APIMethod)
	apiResponseJSON, _ := json.Marshal(apiResponse)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write(apiResponseJSON)
}

type InternalServerErrorHandler struct {
	APIMethod string
}

func (h InternalServerErrorHandler) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	apiResponse := responses.InternalServerError(h.APIMethod)
	apiResponseJSON, _ := json.Marshal(apiResponse)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write(apiResponseJSON)
}

type CustomErrorHandler struct {
	APIMethod string
	Error     error
}

func (h CustomErrorHandler) ServeHTTP(rw http.ResponseWriter, rr *http.Request) {
	apiResponse := responses.NewErrorAPIResponse(h.APIMethod, h.Error)
	apiResponseJSON, err := json.Marshal(apiResponse)
	if err != nil {
		InternalServerErrorHandler{h.APIMethod}.ServeHTTP(rw, rr)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write(apiResponseJSON)
}
