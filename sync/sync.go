package sync

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rebel-l/go-utils/osutils"
)

var (
	ErrFileInfo        = errors.New("failed to read destinations file info")
	ErrCreateDirectory = errors.New("failed to create destination directory")
)

func Do(file File) error {
	destPath, _ := filepath.Split(file.Destination)
	if err := osutils.CreateDirectoryIfNotExists(destPath); err != nil {
		return fmt.Errorf("%w: %v", ErrCreateDirectory, err)
	}

	if err := osutils.CopyFile(file.Source.GetName(), file.Destination); err != nil {
		return err
	}

	return os.Chtimes(file.Destination, file.Source.Info.ModTime(), file.Source.Info.ModTime())
}
