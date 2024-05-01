package storage

import (
	"log"
)

// Interface for memory and database storages
type Repository interface {
	Get(string) (Metric, error)
	Save(Metric) error
}

var storage *MemRepo

func init() {
	log.Println("Init storage")
	storage = NewMemRepo()
}

func Get(metricType MetricType, metricName string) (Metric, bool) {
	pk := pkey(metricType, metricName)
	return storage.Get(pk)
}
