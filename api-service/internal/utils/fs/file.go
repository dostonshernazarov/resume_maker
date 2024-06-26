package fs

import (
	"errors"
	"fmt"
	"os"
)

// EnsureNonEmptyFile ensures file exits and is non empty
func EnsureNonEmptyFile(filePath string) error {
	if filePath == "" {
		return errors.New("missing file path")
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to read file information %v", err))
	}

	if fileInfo.Size() == 0 {
		return errors.New("file is empty")
	}

	return nil
}

// WriteFile writes data to the specified destination.
func WriteFile(destination string, data []byte) error {
	if err := os.WriteFile(destination, data, os.ModePerm); err != nil {
		return errors.New(fmt.Sprintf("failed to write file %v", err))
	}

	return nil
}

// CreateFile creates a file to the specified destination.
func CreateFile(destination string) (*os.File, error) {
	htmlOut, err := os.Create(destination)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create file %v", err))
	}

	return htmlOut, nil
}

// ReadFile reads data from specified file.
func ReadFile(file string) ([]byte, error) {
	fileData, err := os.ReadFile(file)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to read file %v", err))
	}

	return fileData, nil
}
