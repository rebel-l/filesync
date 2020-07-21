package mp3files

import (
	"os"
)

type File struct {
	Name string
	Info os.FileInfo
}
