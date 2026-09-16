package storage

import (
	"errors"
	"fmt"

	interfaces "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/interfaces"
	memorystorage "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/evr-gh/otus-go-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

const (
	GoMemoryStorage = "memory"
	PostgresStorage = "pgx"
)

var ErrUnknownStorageType = errors.New("неизвестный тип хранилища")

func New(storageType string, dsn string) (interfaces.Storage, error) {
	switch storageType {
	case GoMemoryStorage:
		return memorystorage.New(), nil
	case PostgresStorage:
		return sqlstorage.New(storageType, dsn), nil
	default:
		return nil, fmt.Errorf("не создано хранилище с типом %q: %w", storageType, ErrUnknownStorageType)
	}
}
