package storage

import (
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type ImageManager struct {
	baseDir string
}

func NewImageManager(baseDir string) *ImageManager {
	return &ImageManager{
		baseDir: baseDir,
	}
}

func (i *ImageManager) Save(subDir, originalName string, r io.Reader) (string, error) {
	ext := filepath.Ext(originalName)
	name := uuid.NewString() + ext
	rel := filepath.Join(subDir, name)
	full := filepath.Join(i.baseDir, rel)

	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return "", err
	}

	file, err := os.Create(full)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(file, r); err != nil {
		return "", err
	}
	return rel, nil
}

func (i *ImageManager) Delete(key string) error {
	if key == "" {
		return nil
	}

	full := filepath.Join(i.baseDir, key)
	if err := os.Remove(full); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
