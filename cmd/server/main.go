package main

import (
	"flag"
	"fmt"
	"log"

	"metrics/internal/core/service"
	"metrics/internal/infra/api/rest"
	"metrics/internal/infra/store"
)

// func init() {
// flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
// }

func main() {
	var flagRunAddr string

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if err := run(flagRunAddr); err != nil {
		log.Fatal(err)
	}
}

func run(flagRunAddr string) error {
	store, err := store.NewStore()
	if err != nil {
		return fmt.Errorf("failed to initialize a store: %w", err)
	}
	service := service.NewMetricService(store)
	log.Println("Service initialized")
	api := rest.NewAPI(service)

	log.Println("Start Server at:", flagRunAddr)
	if err := api.Run(flagRunAddr); err != nil {
		return fmt.Errorf("server has failed: %w", err)
	}
	return nil
}
