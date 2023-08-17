package mp3files

import (
	"errors"
	"fmt"
	"github.com/rebel-l/mp3sync/config"
	"os"
	"path/filepath"
)

const (
	Extension = ".mp3"
)

var ErrFileList = errors.New("failed to read file list")

func GetFileList(path string, whitelist config.File, blackList config.File) (Files, error) {
	var list Files

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("%w: %v", ErrFileList, err)
		}

		if (len(whitelist) > 0 && !whitelist.Contains(info)) || (len(blackList) > 0 && blackList.Contains(info)) {
			return nil
		}

		list = append(list, File{Name: path, Info: info})

		return nil
	})
	if err != nil {
		return list, err
	}

	return list, nil
}
