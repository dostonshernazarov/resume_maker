package fs

import (
	"errors"
	"fmt"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/logger"
	"os"
)

// EnsureDir ensures directory exits and creates it
func EnsureDir(directory string) {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err := os.Mkdir(directory, os.ModePerm)
		if err != nil {
			logger.Error(errors.New(fmt.Sprintf("EnsureDir - failed to create directory: %v", err)))
		}
	}
}
