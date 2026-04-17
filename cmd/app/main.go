package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sample-web-http/internal/handler"
	redisclient "sample-web-http/internal/redis"
	"sample-web-http/internal/route"
	"sample-web-http/internal/todo"
	"syscall"
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
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Println("Error server:", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("Server shutdown complete")
	rdb.Close()
	log.Println("Redis connection closed")
}
