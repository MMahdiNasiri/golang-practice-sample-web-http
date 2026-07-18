package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sample-web-http/internal/authenticate"
	todoHandler "sample-web-http/internal/handler/todo"
	userHandler "sample-web-http/internal/handler/user"
	"sample-web-http/internal/middleware"
	"sample-web-http/internal/storage/cached"
	"sample-web-http/internal/user"

	postgresclient "sample-web-http/internal/postgres"
	redisclient "sample-web-http/internal/redis"
	"sample-web-http/internal/route"
	postgresRepo "sample-web-http/internal/storage/postgres"
	"sample-web-http/internal/todo"
	"syscall"
	"time"
)

func main() {

	rdb, err := redisclient.New()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Println("Error closing Redis:", err)
		}
	}()

	portgresdb := postgresclient.New()
	defer func(portgresdb *sql.DB) {
		err := portgresdb.Close()
		if err != nil {
			log.Println("Error closing Postgres:", err)
		}
	}(portgresdb)
	todoRepo := postgresRepo.NewTodoRepo(portgresdb)
	cachedRepo := cached.NewCachedTodoRepo(todoRepo, rdb, 360*time.Second)
	todoService := todo.NewService(cachedRepo)
	userRepo := postgresRepo.NewUserRepo(portgresdb)
	userService := user.NewService(userRepo)
	authService := authenticate.NewService()
	authMiddleware := middleware.Auth(authService)

	todoh := &todoHandler.Handler{
		TodoService:    todoService,
		UserService:    userService,
		AuthMiddleware: authMiddleware,
	}
	userh := &userHandler.Handler{
		UserService: userService,
		AuthService: authService,
	}

	router := route.NewRouter(todoh, userh)

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
	err = server.Shutdown(ctx)
	if err != nil {
		log.Println("Server shutdown error:", err)
		return
	}
	log.Println("Server shutdown complete")
}
