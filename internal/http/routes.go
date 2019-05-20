package http

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"github.com/vaksi/messaging/configs"
	"github.com/vaksi/messaging/internal/http/handler"
	"github.com/vaksi/messaging/pkg/logrus"
	"net/http"
)

type Routes struct {
	Config     *configs.Config
	MsgHandler handler.IMessageHandler
}

// Main Router
func (r *Routes) NewRoutes() http.Handler {
	// define route
	router := mux.NewRouter().StrictSlash(false)
	route := router.PathPrefix(r.Config.Api.Prefix).Subrouter()

	// health-check
	route.HandleFunc("/health-check", handler.GetHealthCheck).Methods(http.MethodGet)
	// messages
	route.HandleFunc("/messages", r.MsgHandler.CreateMessage).Methods(http.MethodPost)
	route.HandleFunc("/messages", r.MsgHandler.GetMessageRealtime).Methods(http.MethodGet)
	route.HandleFunc("/messages/inbox", r.MsgHandler.GetMessageInbox).Methods(http.MethodGet)

	// Use Negroni Log Router
	n := negroni.Classic()
	recovery := negroni.NewRecovery() // Panic handler
	if r.Config.App.Debug == false {
		recovery.PrintStack = false
	}

	n.Use(logrus.NewLoggerMiddleware(r.Config.App.Name))
	n.UseHandler(router)
	return n
}
