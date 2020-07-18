package sync

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/rebel-l/go-utils/osutils"
	"github.com/rebel-l/mp3sync/mp3files"
)

var (
	ErrFileInfo        = errors.New("failed to read destinations file info")
	ErrCreateDirectory = errors.New("failed to create destination directory")
)

func Do(source mp3files.File, destination string, preview bool) error {
	if osutils.FileOrPathExists(destination) {
		destInfo, err := os.Lstat(destination)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrFileInfo, err)
		}

		if source.IsInSync(destInfo) {
			return nil
		}
	}

	destPath, _ := filepath.Split(destination)
	if err := osutils.CreateDirectoryIfNotExists(destPath); err != nil {
		return fmt.Errorf("%w: %v", ErrCreateDirectory, err)
	}

	if preview {
		listFormat := color.New(color.FgHiBlue)
		_, _ = listFormat.Println(destination)

		return nil
	}

	return osutils.CopyFile(source.GetName(), destination)
}
