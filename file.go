package main

import (
	"fmt"
	"github.com/rebel-l/go-utils/osutils"
	"time"

	"github.com/rebel-l/mp3sync/filter"
	"github.com/rebel-l/mp3sync/mp3files"
)

func readFileList(path string, confFilter filter.BlackWhiteList) (mp3files.Files, error) {
	if !osutils.FileOrPathExists(path) {
		return nil, fmt.Errorf("%w: %s", errPathNotExisting, path)
	}

	_, _ = description.Print("Read File List: ")
	start := time.Now()

	defer fmt.Println()

	wl := make(filter.File)
	bl := make(filter.File)
	if confFilter != nil {
		wl, _ = confFilter.File(filter.KeyWhitelist)
		bl, _ = confFilter.File(filter.KeyBlacklist)
	}

	fileList, err := mp3files.GetFileList(path, wl, bl)
	if err != nil {
		return nil, err
	}

	duration(start, time.Now(), fmt.Sprintf("%d files found", len(fileList)))

	return fileList, nil
}
