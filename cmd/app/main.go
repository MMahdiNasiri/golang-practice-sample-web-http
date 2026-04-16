package main

import (
	"fmt"
	"net/http"
	"sample-web-http/internal/handler"
	redisclient "sample-web-http/internal/redis"
	"sample-web-http/internal/route"
	"sample-web-http/internal/todo"
)

func main() {
	fmt.Printf("Hello and welcome!\n")
	rdb := redisclient.New()
	todoRepo := todo.NewRepository(rdb)
	todoService := todo.NewService(todoRepo)
	h := &handler.Handler{TodoService: todoService}

	router := route.NewRouter(h)

	http.ListenAndServe(":8080", router)
}
