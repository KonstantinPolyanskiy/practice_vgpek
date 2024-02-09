package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"practice_vgpek/internal/handler"
	"practice_vgpek/internal/repository"
	"practice_vgpek/internal/service"
	"syscall"
	"time"
)

func main() {
	mainCtx, cancel := context.WithCancel(context.Background())

	repo := repository.New()
	services := service.New(repo)
	handlers := handler.New(services)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: handlers.Init(),
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c

		shutdownCtx, _ := context.WithTimeout(mainCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()

			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatalf("graceful shutdown timed out, force exit")
			}
		}()

		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		cancel()
	}()

	err := httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	<-mainCtx.Done()

}
