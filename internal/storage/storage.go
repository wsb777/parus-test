package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

type LocalDiskStorage struct {
	BasePath string
}

func (s *LocalDiskStorage) Save(ctx context.Context, path string, data []byte) error {
	full := filepath.Join(s.BasePath, path)

	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return err
	}

	return os.WriteFile(full, data, 0644)
}

func (s *LocalDiskStorage) Delete(ctx context.Context, path string) error {
	return os.Remove(filepath.Join(s.BasePath, path))
}

func (s *LocalDiskStorage) Read(ctx context.Context, path string) (io.ReadCloser, error) {
	full := filepath.Join(s.BasePath, path)

	return os.Open(full)
}
