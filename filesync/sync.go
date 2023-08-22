package filesync

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rebel-l/go-utils/osutils"
)

var ErrCreateDirectory = errors.New("failed to create destination directory")

func Do(file File) error {
	destPath, _ := filepath.Split(file.Destination.Name)
	if err := osutils.CreateDirectoryIfNotExists(destPath); err != nil {
		return fmt.Errorf("%w: %v", ErrCreateDirectory, err)
	}

	if err := osutils.CopyFile(file.Source.Name, file.Destination.Name); err != nil {
		return err
	}

	return os.Chtimes(file.Destination.Name, file.Source.Info.ModTime(), file.Source.Info.ModTime())
}
