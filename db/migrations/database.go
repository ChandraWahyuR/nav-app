package migrations

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"
)

func CreateTables(db *sql.DB) error {
	files := []string{
		"./db/migrations/001_UserModel.sql",
		"./db/migrations/002_OtpModel.sql",
	}

	for _, v := range files {
		query, err := os.ReadFile(v)
		if err != nil {
			log.Printf("Gagal membaca file %s: %s", v, err)
			continue
		}

		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()

		res, err := db.ExecContext(ctx, string(query))
		if err != nil {
			log.Printf("Error %s when creating table", err)
			continue
		}
		rows, err := res.RowsAffected()
		if err != nil {
			log.Printf("Error %s when getting rows affected", err)
			continue
		}
		log.Printf("Rows affected when creating table: %d", rows)
	}
	return nil
}
