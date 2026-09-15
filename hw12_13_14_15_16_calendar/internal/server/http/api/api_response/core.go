package api

import (
	"encoding/json"
	"fmt"
)

type DataItem struct {
	Item any `json:"item"`
}
type DataItems struct {
	Items any `json:"items"`
}

type Response struct {
	APIMethod string `json:"method"`
	Error     error  `json:"error"`
	Data      any    `json:"data"`
}

func (r Response) MarshalJSON() ([]byte, error) {
	var errorMsg string
	if r.Error != nil {
		errorMsg = r.Error.Error()
	}
	anon := struct {
		APIMethod string `json:"method"`
		Error     string `json:"error"`
		Data      any    `json:"data"`
	}{
		APIMethod: r.APIMethod,
		Error:     errorMsg,
		Data:      r.Data,
	}
	return json.Marshal(anon)
}

func NewAPIResponse(apiMethod string) *Response {
	r := new(Response)
	r.APIMethod = apiMethod
	return r
}

func NewErrorAPIResponse(apiMethod string, err error) *Response {
	r := new(Response)
	r.APIMethod = apiMethod
	r.Error = err
	return r
}

func NewErrorAPIResponseByString(apiMethod string, format string, a ...any) *Response {
	r := new(Response)
	r.APIMethod = apiMethod
	r.Error = fmt.Errorf(format, a...)
	return r
}

func InternalServerError(apiMethod string) *Response {
	r := new(Response)
	r.APIMethod = apiMethod
	r.Error = fmt.Errorf("internal server error")
	return r
}

func InvalidHTTPMethod(apiMethod string) *Response {
	r := new(Response)
	r.APIMethod = apiMethod
	r.Error = fmt.Errorf("invalid HTTP method")
	return r
}

func InvalidRequestBody(apiMethod string) *Response {
	r := new(Response)
	r.APIMethod = apiMethod
	r.Error = fmt.Errorf("invalid HTTP request body")
	return r
}
