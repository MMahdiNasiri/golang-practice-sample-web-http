package route

import (
	"net/http"
	"sample-web-http/internal/handler"
	"sample-web-http/internal/middleware"
)

func NewRouter(h *handler.Handler) *http.ServeMux {
	server := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	server.Handle("/static/", http.StripPrefix("/static/", fs))

	server.HandleFunc("/", middleware.RateLimit(h.PageHandler))
	server.HandleFunc("/all/", h.ListContext)
	server.HandleFunc("/create", h.CreateContext)
	server.HandleFunc("/update", h.UpdateContext)
	server.HandleFunc("/delete", h.DeleteContext)

	return server
}
