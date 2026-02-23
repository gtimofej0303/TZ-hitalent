package main

import (
	"context"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gtimofej0303/TZ-hitalent/config"
	"github.com/gtimofej0303/TZ-hitalent/internal/domain"
	"github.com/gtimofej0303/TZ-hitalent/internal/repository/mygorm"
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

	goose.SetDialect("postgres")
	if err := goose.Up(gormSqlDB, "./migrations"); err != nil {
        log.Fatal("goose up failed:", err)
    }
    log.Println("Migrations applied succesfully")

	i := 1
	for i == 1 {
		continue
	}
}
