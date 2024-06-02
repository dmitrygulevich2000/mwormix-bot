package storage

import (
	"io"
	"os"
	"time"
)

type SearchesStorage struct {
	file   *os.File
	cached time.Time
}

func NewSearches(filename string) (SearchesStorage, error) {
	// rw- r-- ---
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		return SearchesStorage{}, err
	}
	storage := SearchesStorage{
		file: file,
	}
	return storage, nil
}

func (s *SearchesStorage) Close() {
	_ = s.file.Close()
}

func (s *SearchesStorage) StoreSearch(at time.Time) error {
	err := s.dump(at)
	if err != nil {
		return err
	}
	s.cached = at
	return nil
}

func (s *SearchesStorage) LastSearchTime() (time.Time, error) {
	if !s.cached.IsZero() {
		return s.cached, nil
	}
	err := s.load()
	return s.cached, err
}

func (s *SearchesStorage) dump(ts time.Time) error {
	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	err = s.file.Truncate(0)
	if err != nil {
		return err
	}
	data, err := ts.MarshalText()
	if err != nil {
		return err
	}
	_, err = s.file.Write(data)
	return err
}

func (s *SearchesStorage) load() error {
	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(s.file)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	err = s.cached.UnmarshalText(data)
	return err
}
