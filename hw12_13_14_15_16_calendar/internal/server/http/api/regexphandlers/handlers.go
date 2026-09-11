package regexhandlers

import (
	"net/http"

	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
)

type RegexpHandler struct {
	qpp     QueryPathPattern
	handler http.Handler
}

func NewRegexpHandler(pattern string, params Params, handler http.Handler) *RegexpHandler {
	rh := new(RegexpHandler)
	rh.qpp = *NewQueryPathPattern(pattern, params)
	rh.handler = handler
	return rh
}

type RegexpHandlers struct {
	defaultHandler http.Handler
	logger         interfaces.Logger
	app            interfaces.Application
	crossroad      []RegexpHandler
}

func NewRegexpHandlers(h http.Handler, log interfaces.Logger, app interfaces.Application,
	rh ...RegexpHandler,
) RegexpHandlers {
	regexpHandlers := new(RegexpHandlers)
	regexpHandlers.defaultHandler = h
	regexpHandlers.logger = log
	regexpHandlers.app = app
	regexpHandlers.crossroad = rh
	return *regexpHandlers
}

func (rhs RegexpHandlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlerWasNotFound := true
	for _, rh := range rhs.crossroad {
		if rh.qpp.match(r.URL.Path) {
			handlerWasNotFound = false
			r.Form = rh.qpp.GetValues(r.URL.Path)
			rh.handler.ServeHTTP(w, r)
			break
		}
	}
	if handlerWasNotFound && rhs.defaultHandler != nil {
		rhs.defaultHandler.ServeHTTP(w, r)
	}
}
