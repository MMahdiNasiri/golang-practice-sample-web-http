package main

import (
	"net/http"
	"sample-web-http/internal/handler"
	redisclient "sample-web-http/internal/redis"
	"sample-web-http/internal/route"
	"sample-web-http/internal/todo"
	"time"
)

func main() {

	rdb := redisclient.New()
	todoRepo := todo.NewRepository(rdb)
	todoService := todo.NewService(todoRepo)
	h := &handler.Handler{TodoService: todoService}

	router := route.NewRouter(h)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	server.ListenAndServe()
}
