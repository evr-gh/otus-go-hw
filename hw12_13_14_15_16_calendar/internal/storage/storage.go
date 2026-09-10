package storage

import (
	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
	memorystorage "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var (
	MemoryStorage   = "memory"
	PostgresStorage = "pgx"
)

func New(storageType string, dsn string) interfaces.Storage {
	switch storageType {
	case MemoryStorage:
		return memorystorage.New()
	default:
		return sqlstorage.New(storageType, dsn)
	}
}
