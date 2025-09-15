package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gfrei/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	db, _ := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	server := newServer(dbQueries, platform)
	log.Fatal(server.ListenAndServe())
}
