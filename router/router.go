package router

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Path        string
	StrictSlash bool
	Handler     http.Handler
}

func New(routes []Route) *mux.Router {
	router := mux.NewRouter()

	for _, route := range routes {
		// Register defined route
		router.
			StrictSlash(route.StrictSlash).
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(route.Handler)

		// Register route with method OPTIONS
		router.
			StrictSlash(route.StrictSlash).
			Methods(http.MethodOptions).
			Path(route.Path).
			Name(route.Name).
			Handler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
	}

	makeSwaggerHandler(router)

	return router
}

func makeSwaggerHandler(r *mux.Router) {
	const docsPath = "/docs"

	r.StrictSlash(false).Path(docsPath).Handler(http.RedirectHandler(docsPath+"/", http.StatusMovedPermanently))

	r.StrictSlash(true).PathPrefix(docsPath + "/").Handler(
		http.StripPrefix(docsPath+"/", http.FileServer(http.Dir("./swagger"))),
	)

	r.Path("/api-docs").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./swagger.yaml")
	})
}
