package main

import (
	"database/sql"
	"log"

	"github.com/pressly/goose"
	_ "github.com/lib/pq"

	"github.com/gtimofej0303/TZ-hitalent/config"
)

func main() {
	cfg := config.Load()
	if cfg == nil {
		log.Fatal("config is empty")
	}

	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		log.Fatal("failed to open db: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect to db: ", err)
	}

	goose.SetDialect("postgres")
	if err := goose.Up(db, "./migrations"); err != nil {
		log.Fatal("failed to run migrations: ", err)
	}

	log.Println("Migrations applied successfully")

	i := 1
	for i == 1 {
		continue
	}
}
