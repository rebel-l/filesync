package sync

import (
	"fmt"
	"os"
	"time"

	"github.com/rebel-l/go-utils/osutils"
	"github.com/rebel-l/mp3sync/mp3files"
)

type File struct {
	Source      mp3files.File
	Destination string
}

func (f File) IsInSync() (bool, error) {
	if !osutils.FileOrPathExists(f.Destination) {
		return false, nil
	}

	destInfo, err := os.Lstat(f.Destination)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrFileInfo, err)
	}

	if f.Source.Info.Size() != destInfo.Size() {
		return false, nil
	}

	if !timeEqual(f.Source.Info.ModTime(), destInfo.ModTime()) {
		return false, nil
	}

	return true, nil
}

func timeEqual(a, b time.Time) bool {
	d := a.Sub(b)
	return d.Seconds() > -5 && d.Seconds() < 5
}
