package filesync

import (
	"time"

	"github.com/rebel-l/mp3sync/mp3files"
)

type File struct {
	Source      mp3files.File
	Destination mp3files.File
}

func (f File) IsInSync() bool {
	if f.Source.Info == nil || f.Destination.Info == nil {
		return false
	}

	if f.Source.Info.Size() != f.Destination.Info.Size() {
		return false
	}

	if !timeEqual(f.Source.Info.ModTime(), f.Destination.Info.ModTime()) {
		return false
	}

	return true
}

func timeEqual(a, b time.Time) bool {
	d := a.Sub(b)
	return d.Seconds() > -5 && d.Seconds() < 5
}
