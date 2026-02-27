package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gtimofej0303/TZ-hitalent/config"
	"github.com/gtimofej0303/TZ-hitalent/internal/handler"
)

func main() {
	cfg := config.Load()
	if cfg == nil {
		log.Fatal("config is empty")
	}

	gormDB, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to open db: ", err)
	}

	gormSqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal("failed to connect to DB: ", err)
	}

	if err := gormSqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	const migrationsDir = "./migrations"

	goose.SetDialect("postgres")

	//Для тестирования
	/*if err := goose.Reset(gormSqlDB, migrationsDir); err != nil {
        log.Fatal("goose reset failed: ", err)
    }*/

	if err := goose.Up(gormSqlDB, "./migrations"); err != nil {
		log.Fatal("goose up failed:", err)
	}
	log.Println("Migrations applied successfully")

	port := "8080"

	router := handler.NewRouter(gormDB)

	log.Printf("Server starting on :%s", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
        log.Fatal("server error: ", err)
    }
}
