package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FileStore struct {
	base string
}

func NewFileStore(base string) *FileStore {
	_ = os.MkdirAll(base, 0o755)
	return &FileStore{base: base}
}

func (s *FileStore) reqDir(id string) string {
	return filepath.Join(s.base, id)
}
func (s *FileStore) verPath(id string, version uint64) string {
	return filepath.Join(s.reqDir(id), fmt.Sprintf("v%d.sexpr", version))
}
func (s *FileStore) latestPath(id string) string {
	return filepath.Join(s.reqDir(id), "latest")
}

func (s *FileStore) Put(id string, version uint64, text string) error {
	if err := os.MkdirAll(s.reqDir(id), 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	if err := os.WriteFile(s.verPath(id, version), []byte(text), 0o644); err != nil {
		return fmt.Errorf("failed to write version file: %w", err)
	}
	if err := os.WriteFile(s.latestPath(id), []byte(fmt.Sprintf("%d", version)), 0o644); err != nil {
		return fmt.Errorf("failed to write latest file: %w", err)
	}
	return nil
}

func (s *FileStore) GetLatest(id string) (uint64, string, error) {
	b, err := os.ReadFile(s.latestPath(id))
	if err != nil {
		return 0, "", err
	}
	v, err := strconv.ParseUint(strings.TrimSpace(string(b)), 10, 64)
	if err != nil {
		return 0, "", err
	}
	txt, err := os.ReadFile(s.verPath(id, v))
	if err != nil {
		return 0, "", err
	}
	return v, string(txt), nil
}

func (s *FileStore) Get(id string, version uint64) (string, error) {
	b, err := os.ReadFile(s.verPath(id, version))
	if err != nil {
		return "", err
	}
	return string(b), nil
}
