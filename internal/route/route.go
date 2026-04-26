package route

import (
	"net/http"
)

type RouteRegistrar interface {
	RegisterRoutes(mux *http.ServeMux)
}

func NewRouter(registrars ...RouteRegistrar) *http.ServeMux {
	server := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	server.Handle("/static/", http.StripPrefix("/static/", fs))

	for _, r := range registrars {
		r.RegisterRoutes(server)
	}

	return server
}
