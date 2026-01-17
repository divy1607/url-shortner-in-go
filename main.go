package main

import (
	"log"
	"net/http"
)

func main() {
	initDB()
	initRedis()

	http.HandleFunc("/shorten", createShortURL)
	http.HandleFunc("/", redirect)

	log.Println("server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
