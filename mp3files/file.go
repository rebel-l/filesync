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
