package main

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"practice_vgpek/internal/dao"
	"practice_vgpek/internal/handler"
	"practice_vgpek/internal/service"
	"practice_vgpek/pkg/logger"
	"practice_vgpek/pkg/postgres"
	"syscall"
	"time"
)

// @title						ВГПЭК API
// @version					0.1
// @description				API для работы с практическими заданиями
// @host						localhost:8080
// @basePath					/
// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	mainCtx, cancel := context.WithCancel(context.Background())

	err := mustLoadConfig()
	if err != nil {
		panic(err)
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
	if err != nil {
		logging.Fatal("error connect to db", zap.Error(err))
	}

	err = migrateDB(db)
	if err != nil {
		logging.Error("ошибка миграции", zap.Error(err))
	}

	dao := dao.New(db, logging)
	services := service.New(dao, logging)
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

func migrateDB(pool *pgxpool.Pool) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(pool)
	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}
