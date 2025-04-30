package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func InitDatabase(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Database.dbUser, cfg.Database.dbPass, cfg.Database.dbHost, cfg.Database.dbPort, cfg.Database.dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfull connect to database on DSN:%s", dsn)
	return db, nil
}
