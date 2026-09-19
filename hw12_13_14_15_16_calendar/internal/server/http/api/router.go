package api

import (
	"net/http"

	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
	commonHandlers "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/api/handlers/common"
	apiRequests "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/api/handlers/requests"
	regexped "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/api/regexphandlers"
	middleware "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/server/http/middleware"
)

type DefaultHandler struct{}

func (h DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/api/version", http.StatusTemporaryRedirect)
}

var (
	none = regexped.Params{}
	id   = regexped.Params{"id"}
)

func Handlers(logger interfaces.Logger, app interfaces.Application) regexped.RegexpHandlers {
	return regexped.NewRegexpHandlers(
		middleware.Instance().Listen(DefaultHandler{}),
		logger,
		app,
		*regexped.NewRegexpHandler(
			`/api/version`,
			none,
			middleware.Instance().Listen(commonHandlers.VersionHandler{}),
		),
		*regexped.NewRegexpHandler(
			`/api/events/`,
			none,
			middleware.Instance().Listen(apiRequests.EventsListHandler{
				Logger: logger, App: app, APIMethod: "api.events.list",
			}),
		),
		*regexped.NewRegexpHandler(
			`/api/events/notsheduled`,
			none,
			middleware.Instance().Listen(apiRequests.EventsListHandler{
				Logger: logger, App: app, APIMethod: "api.events.listnotsheduled",
			}),
		),
		*regexped.NewRegexpHandler(
			`/api/events/create`,
			none,
			middleware.Instance().Listen(apiRequests.EventsCreateHandler{
				Logger: logger, App: app,
			}),
		),
		*regexped.NewRegexpHandler(
			`/api/events/{numeric}`,
			id,
			middleware.Instance().Listen(apiRequests.EventsGetHandler{
				Logger: logger, App: app,
			}),
		),
		*regexped.NewRegexpHandler(
			`/api/events/{numeric}/update`,
			id,
			middleware.Instance().Listen(apiRequests.EventsUpdateHandler{
				Logger: logger, App: app,
			}),
		),
		*regexped.NewRegexpHandler(
			`/api/events/{numeric}/delete`,
			id,
			middleware.Instance().Listen(apiRequests.EventsDeleteHandler{
				Logger: logger, App: app,
			}),
		),
	)
}
