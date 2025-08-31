package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/davidegagliardi/syowatchdog/pkg/types"
)

type FileStorage struct {
	basePath string
	filename string
}

func NewFileStorage(basePath string) *FileStorage {
	return &FileStorage{
		basePath: basePath,
		filename: "image_data.json",
	}
}

func (fs *FileStorage) ensureDirectory() error {
	return os.MkdirAll(fs.basePath, 0755)
}

func (fs *FileStorage) getFilePath() string {
	return filepath.Join(fs.basePath, fs.filename)
}

func (fs *FileStorage) SaveImageData(data *types.ImageData) error {
	if err := fs.ensureDirectory(); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	filePath := fs.getFilePath()
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (fs *FileStorage) LoadImageData() (*types.ImageData, error) {
	filePath := fs.getFilePath()

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, nil // File doesn't exist, return nil (first run)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var imageData types.ImageData
	if err := json.Unmarshal(data, &imageData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return &imageData, nil
}

func (fs *FileStorage) DeleteImageData() error {
	filePath := fs.getFilePath()
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
