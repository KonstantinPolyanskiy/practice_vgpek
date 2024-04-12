package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"practice_vgpek/pkg/postgres"
)

const baseApiUrl = "localhost:8080"

func mustLoadConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("local")

	return viper.ReadInConfig()
}

type InitUtils struct {
	baseActions [4]string
	baseRoles   [3]string
	baseObjects [5]string

	db *pgxpool.Pool
}

func NewInitUtils(db *pgxpool.Pool, roles [3]string, actions [4]string, objects [5]string) InitUtils {
	return InitUtils{
		db:          db,
		baseRoles:   roles,
		baseActions: actions,
		baseObjects: objects,
	}
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

	roles := [3]string{"ADMIN", "TEACHER", "STUDENT"}
	actions := [4]string{"ADD", "GET", "EDIT", "DEL"}
	objects := [5]string{"KEY", "RBAC", "ISSUED PRACTICE", "SOLVED PRACTICE", "MARK"}

	iu := NewInitUtils(db, roles, actions, objects)

	iu.createBaseRoles()
	iu.createBaseActions()
	iu.createBaseObjects()

	for i := 1; i <= 5; i++ {
		iu.setAdminPerm(i)
	}

}

func (u InitUtils) createBaseRoles() {
	insertQuery := `INSERT INTO internal_role (role_name) VALUES ($1)`

	for _, role := range u.baseRoles {
		_, err := u.db.Exec(context.Background(), insertQuery, role)
		if err != nil {
			log.Printf("Возникла ошибка %v при вставке роли %s", err, role)
		}
	}

	log.Println("Основные роли успешно вставлены")
}

func (u InitUtils) createBaseActions() {
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

func (u InitUtils) createBaseObjects() {
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

func (u InitUtils) setAdminPerm(objId int) {
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
