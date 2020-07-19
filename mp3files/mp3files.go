package mp3files

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	Extension = ".mp3"
)

var (
	ErrFileList = errors.New("failed to read file list")
)

func GetFileList(path string) (Files, error) {
	var list Files

	i := 0

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("%w: %v", ErrFileList, err)
		}

		if strings.ToLower(filepath.Ext(info.Name())) == Extension {
			i++
			dir, _ := filepath.Split(path)
			list = append(list, File{Path: dir, Info: info})
		}

		if i > 50 { // nolint:gomnd
			// TODO: remove, just temporary
			return fmt.Errorf("100 files reached")
		}

		return nil
	})

	if err != nil {
		return list, err
	}

	return list, nil
}
