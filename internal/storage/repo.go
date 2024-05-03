package storage

import (
	"log"
)

// Interface for memory and database storages
type Repository interface {
	Get(string) (interface{}, error)
	Save(interface{}) error
}

var Storage *MemRepo

func init() {
	log.Println("Init storage")
	Storage = NewMemRepo()
}
