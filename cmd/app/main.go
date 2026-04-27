package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sample-web-http/internal/authenticate"
	todoHandler "sample-web-http/internal/handler/todo"
	userHandler "sample-web-http/internal/handler/user"
	"sample-web-http/internal/middleware"
	"sample-web-http/internal/user"

	//redisclient "sample-web-http/internal/redis"
	"sample-web-http/internal/route"
	//redisRepo "sample-web-http/internal/storage/redis"
	postgresclient "sample-web-http/internal/postgres"
	postgresRepo "sample-web-http/internal/storage/postgres"
	"sample-web-http/internal/todo"
	"syscall"
	"time"
)

func main() {

	//rdb := redisclient.New()
	//todoRepo := redisRepo.NewTodoRepo(rdb)
	portgresdb := postgresclient.New()
	defer portgresdb.Close()
	todoRepo := postgresRepo.NewTodoRepo(portgresdb)
	todoService := todo.NewService(todoRepo)
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
	server.Shutdown(ctx)
	log.Println("Server shutdown complete")
	//rdb.Close()
	log.Println("Redis connection closed")
}
