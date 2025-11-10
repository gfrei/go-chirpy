package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gfrei/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Start server")

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	polkaKey := os.Getenv("POLKA_KEY")

	db, _ := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	server := newServer(dbQueries, platform, secret, polkaKey)
	log.Fatal(server.ListenAndServe())
}
