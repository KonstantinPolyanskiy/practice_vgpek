package main

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"practice_vgpek/pkg/postgres"
)

const baseApiUrl = "127.0.0.1:8080"

func mustLoadConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("local")

	return viper.ReadInConfig()
}

func main() {
	err := mustLoadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}

	db, err := postgres.NewPostgresPool(postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatal("error connect to db", zap.Error(err))
	}
	defer db.Close()

}
