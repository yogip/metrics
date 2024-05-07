package main

import (
	"log"

	"metrics/internal/agent"
)

func main() {
	log.Println("Start agent")
	agent.Run()
}
