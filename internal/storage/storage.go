package storage

import (
	"path"
	"time"
)

type Storage interface {
	Close()
	LastSearchTime() (time.Time, error)
	StoreSearch(at time.Time) error
}

type Config struct {
	SearchesDataFolder string `split_words:"true" required:"true"`
	SearchesDBName     string `split_words:"true" required:"true"`
	// sql db configs
}

func New(searches SearchesStorage) Storage {
	return &storageImpl{
		SearchesStorage: searches,
	}
}

func NewFromConfig(config Config) (Storage, error) {
	searches, err := NewSearches(path.Join(config.SearchesDataFolder, config.SearchesDBName))
	if err != nil {
		return nil, err
	}

	return New(searches), nil
}

type storageImpl struct {
	SearchesStorage
	// sql db connection
}

func (s *storageImpl) Close() {
	s.SearchesStorage.Close()
}
