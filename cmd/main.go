package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"main.go/internal/handler"
	"main.go/internal/service"
	"main.go/internal/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		cancel()
	}()
	Run(ctx)
}

func Run(ctx context.Context) {
	db, err := storage.InitDB()
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("can't close db err: %v\n", err)
		} else {
			log.Printf("db closed")
		}
	}()

	storages := storage.NewStorage(db)
	services := service.NewService(storages)
	handlers := handler.NewHandler(services)
	handlers.InitRouter()

	server := http.Server{
		Addr:         ":8023",
		Handler:      handlers.Mux,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 30,
	}

	fmt.Println("Starting server on http://localhost:8023")
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
			cancel()
			return
		}
	}()
	<-ctx.Done()
	ctx, cancel = context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		log.Println(err)
		return
	}
}
