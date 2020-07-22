package mp3files

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rebel-l/mp3sync/config"
)

const (
	Extension = ".mp3"
)

var ErrFileList = errors.New("failed to read file list")

func GetFileList(path string, filter config.Filter) (Files, error) {
	var list Files

	i := 0

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("%w: %v", ErrFileList, err)
		}

		if strings.ToLower(filepath.Ext(info.Name())) == Extension {
			if filter.Contains(path) {
				return nil
			}

			i++
			list = append(list, File{Name: path, Info: info})
		}

		return nil
	})
	if err != nil {
		return list, err
	}

	return list, nil
}
