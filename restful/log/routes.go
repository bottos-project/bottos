package logrestful

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	SET_CONFIGITEM = "/v1/log/setconfigitem"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	//wallet
	Route{
		"",
		"POST",
		SET_CONFIGITEM,
		SetConfigItem,
	},
}
