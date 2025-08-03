package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/junwei890/crawler/internal/database"
	"github.com/junwei890/crawler/src"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	queries := database.New(db)

	if err := src.Init(queries); err != nil {
		log.Fatal(err)
	}
}
