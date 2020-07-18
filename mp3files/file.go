package mp3files

import (
	"os"
	"path/filepath"
)

type File struct {
	Path string
	Info os.FileInfo
}

func (f File) GetName() string {
	return filepath.Join(f.Path, f.Info.Name())
}

func (f File) IsInSync(destInfo os.FileInfo) bool {
	if f.Info.Size() != destInfo.Size() {
		return false
	}

	if !f.Info.ModTime().Equal(destInfo.ModTime()) {
		return false
	}

	return true
}
