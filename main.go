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
	db, _ := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	server := newServer(dbQueries)
	log.Fatal(server.ListenAndServe())
}
