package main

import (
	"log"
)

func main() {
	server := newServer()

	log.Fatal(server.ListenAndServe())
}
