package storage

import (
	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
	memorystorage "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var (
	GoMemoryStorage = "memory"
	PostgresStorage = "sql"
)

func New(storageType string, a ...any) interfaces.Storage {
	switch storageType {
	case GoMemoryStorage:
		return memorystorage.New()
	case PostgresStorage:
		if len(a) > 0 {
			return sqlstorage.New(a[0].(string))
		}
		return sqlstorage.New("")
	default:
		return memorystorage.New()
	}
}
