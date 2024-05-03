package main

import (
	"log"

	"github.com/yogip/metrics/internal/agent"
)

func main() {
	log.Println("Start agent")
	agent.Run()
}
