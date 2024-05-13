package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"metrics/internal/core/service"
	"metrics/internal/infra/api/rest"
	"metrics/internal/infra/store"
)

func init() {

}

func main() {
	var runAddress string

	flag.StringVar(&runAddress, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if address, ok := os.LookupEnv("ADDRESS"); ok {
		runAddress = address
	}

	if err := run(runAddress); err != nil {
		log.Fatal(err)
	}
}

func run(runAddress string) error {
	store, err := store.NewStore()
	if err != nil {
		return fmt.Errorf("failed to initialize a store: %w", err)
	}
	service := service.NewMetricService(store)
	log.Println("Service initialized")
	api := rest.NewAPI(service)

	log.Println("Start Server at:", runAddress)
	if err := api.Run(runAddress); err != nil {
		return fmt.Errorf("server has failed: %w", err)
	}
	return nil
}
