package server

import (
	"log"
	"net/http"

	"github.com/yogip/metrics/internal/server/handlers"
)

func Run() {
	log.Println("Start server")
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handlers.UpdateHandler)
	mux.HandleFunc("/value/", handlers.GetHandler)

	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
