package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		panic("PORT must be set!")
	}
	port = ":" + port
	mux := http.NewServeMux()
	mux.HandleFunc("/", homePageHandle)

	log.Printf("Starting web server on: %s", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
