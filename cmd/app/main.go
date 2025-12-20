package main

import (
	"log"
	"net/http"

	"github.com/viacheslaev/url-shortener/internal/feature/link"
	"github.com/viacheslaev/url-shortener/internal/server"
)

func main() {
	service := link.NewURLService()
	router := server.NewRouter(service)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
