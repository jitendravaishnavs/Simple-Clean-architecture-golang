package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("DB_DSN")) // e.g., root:pass@tcp(127.0.0.1:3306)/dbname
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	return db
}
