/*
	CS686 Project 3 and 4: Build a Gossip Network and a simple PoW.
	Create a new router.
	Author: Kei Fukutani
	Date  : March 31, Apr 6, 2019
*/
package p3

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Create a new router (HTTP server) and register the handlers.
func NewRouter(routes []Route) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	return router
}
