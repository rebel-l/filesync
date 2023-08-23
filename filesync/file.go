package filesync

import (
	"os"
	"time"
)

const (
	OperationCopy   = "copy"
	OperationDelete = "delete"
)

type File struct {
	Source      os.FileInfo
	Destination os.FileInfo
	Operation   string
}

func (f File) IsInSync() bool {
	if f.Source == nil || f.Destination == nil {
		return false
	}

	if f.Source.Size() != f.Destination.Size() {
		return false
	}

	if !timeEqual(f.Source.ModTime(), f.Destination.ModTime()) { // TODO: ensure destination file has same ModeTime as source file or ignore this
		return false
	}

	return true
}

func timeEqual(a, b time.Time) bool {
	d := a.Sub(b)
	return d.Seconds() > -5 && d.Seconds() < 5
}
