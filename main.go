package main

import "net/http"

func main() {
	serverMux := http.NewServeMux()
	server := http.Server{}

	server.Handler = serverMux
	server.Addr = ":8080"

	server.ListenAndServe()
}
