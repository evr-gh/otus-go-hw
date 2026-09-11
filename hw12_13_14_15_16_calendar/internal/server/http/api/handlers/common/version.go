package api

import (
	"encoding/json"
	"net/http"

	responses "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/api/api_response"
)

type VersionHandler struct{}

func (h VersionHandler) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	apiResponse := responses.NewAPIResponse("api.version")
	apiResponse.Data = struct {
		Version string
	}{
		Version: "1.0.0",
	}
	apiResponseJSON, _ := json.Marshal(apiResponse)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(apiResponseJSON)
}
