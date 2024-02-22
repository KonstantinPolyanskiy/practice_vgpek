package main

import (
	"context"
	"errors"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"practice_vgpek/internal/handler"
	"practice_vgpek/internal/repository"
	"practice_vgpek/internal/service"
	"practice_vgpek/pkg/logger"
	"practice_vgpek/pkg/postgres"
	"syscall"
	"time"
)

func main() {
	mainCtx, cancel := context.WithCancel(context.Background())

	err := mustLoadConfig()
	if err != nil {
		panic("ошибка в чтении конфига")
	}

	logCfg := logger.Config{
		Level:         "debug",
		HasCaller:     true,
		HasStacktrace: true,
		Encoding:      "json",
	}

	logging, _ := logger.New(logCfg)
	defer logging.Sync()

	db, err := postgres.NewPostgresPool(postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	repo := repository.New(db)
	services := service.New(repo)
	handlers := handler.New(services, logging)

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

	err = httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	<-mainCtx.Done()

}

func mustLoadConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("local")

	return viper.ReadInConfig()
}
