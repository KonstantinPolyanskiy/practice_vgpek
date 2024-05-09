package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
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
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/service"
	"practice_vgpek/pkg/logger"
	"practice_vgpek/pkg/postgres"
	"practice_vgpek/pkg/rndutils"
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

	initBase(db)

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

func initBase(db *pgxpool.Pool) {
	roles := [3]string{"ADMIN", "TEACHER", "STUDENT"}
	actions := [4]string{"ADD", "GET", "EDIT", "DEL"}
	objects := [5]string{"KEY", "RBAC", "ISSUED PRACTICE", "SOLVED PRACTICE", "MARK"}

	iu := newInitUtils(db, roles, actions, objects)

	iu.createBaseRoles()
	iu.createBaseActions()
	iu.createBaseObjects()

	for i := 1; i <= 5; i++ {
		iu.setAdminPerm(i)
	}

	iu.createAdminKey(1)
}

type initUtils struct {
	baseActions [4]string
	baseRoles   [3]string
	baseObjects [5]string

	db *pgxpool.Pool
}

func newInitUtils(db *pgxpool.Pool, roles [3]string, actions [4]string, objects [5]string) initUtils {
	return initUtils{
		db:          db,
		baseRoles:   roles,
		baseActions: actions,
		baseObjects: objects,
	}
}

func (u initUtils) createBaseRoles() {
	insertQuery := `INSERT INTO internal_role (role_name) VALUES ($1)`

	for _, role := range u.baseRoles {
		_, err := u.db.Exec(context.Background(), insertQuery, role)
		if err != nil {
			log.Printf("Возникла ошибка %v при вставке роли %s", err, role)
		}
	}

	log.Println("Основные роли успешно вставлены")
}

func (u initUtils) createBaseActions() {
	insertQuery := `
	INSERT INTO internal_action (internal_action_name)
	VALUES ($1)`

	for _, action := range u.baseActions {
		_, err := u.db.Exec(context.Background(), insertQuery, action)
		if err != nil {
			log.Printf("Возникла ошибка %v при вставке действия %s", err, action)
		}
	}

	log.Println("Основные действия успешно вставлены")
}

func (u initUtils) createBaseObjects() {
	insertQuery := `
	INSERT INTO internal_object (internal_object_name)
	VALUES ($1)`

	for _, object := range u.baseObjects {
		_, err := u.db.Exec(context.Background(), insertQuery, object)
		if err != nil {
			log.Printf("Возникла ошибка %v при вставке объекта %s", err, object)
		}
	}

	log.Println("Основные объекты успешно вставлены")
}

func (u initUtils) setAdminPerm(objId int) {
	insertPermQuery := `
	INSERT INTO role_permission (internal_role_id, internal_action_id, internal_object_id) 
	VALUES ($1, $2, $3) 
`
	for i, _ := range u.baseActions {
		i += 1
		_, err := u.db.Exec(context.Background(), insertPermQuery, 1, i, objId)
		if err != nil {
			log.Printf("Возникла ошибка %v при установки доступа", err)
		}
	}

	log.Println("Доступы администратора успешно добавлены")
}

func (u initUtils) createAdminKey(usages int) {
	var insertedKeyId int

	insertKeyQuery := `
	INSERT INTO registration_key (internal_role_id, body_key, max_count_usages, current_count_usages, created_at)  
	VALUES (@RoleId, @BodyKey, @MaxCountUsages, @CurrentCountUsages, @CreatedAt)
	RETURNING reg_key_id
	`

	args := pgx.NamedArgs{
		"RoleId":             1,
		"BodyKey":            rndutils.RandNumberString(5) + rndutils.RandString(5),
		"MaxCountUsages":     usages,
		"CurrentCountUsages": 0,
		"CreatedAt":          time.Now(),
	}

	err := u.db.QueryRow(context.Background(), insertKeyQuery, args).Scan(&insertedKeyId)
	if err != nil {
		log.Printf("Возникла ошибка %v при создании ключа доступа", err)
	}

	getKeyQuery := `SELECT * FROM registration_key WHERE reg_key_id = $1`

	row, err := u.db.Query(context.Background(), getKeyQuery, insertedKeyId)
	if err != nil {
		log.Printf("Возникла ошибка %v при создании ключа доступа", err)
	}

	savedKey, err := pgx.CollectOneRow(row, pgx.RowToStructByName[entity.Key])
	if err != nil {
		log.Printf("Возникла ошибка %v при создании ключа доступа", err)
	}

	file, err := os.Create("key.json")
	if err != nil {
		log.Printf("Ошибка создания файла - %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(savedKey)
	if err != nil {
		log.Printf("Ошибка кодирования структуры - %v", err)
		return
	}
}
