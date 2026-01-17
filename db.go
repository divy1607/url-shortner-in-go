package main

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error

	godotenv.Load()

	db, err = sql.Open(
		"postgres",
		os.Getenv("DATABASE_URL"),
	)

	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

}
