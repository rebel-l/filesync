package main

import (
	"fmt"
	"time"

	"github.com/rebel-l/go-utils/osutils"
	"github.com/rebel-l/mp3sync/filter"
	"github.com/rebel-l/mp3sync/mp3files"
)

type FileChannel struct {
	Files mp3files.Files
	Err   error
}

func readFileList(res chan FileChannel, path string, confFilter filter.BlackWhiteList) {
	if !osutils.FileOrPathExists(path) {
		res <- FileChannel{Err: fmt.Errorf("%w: %s", errPathNotExisting, path)}
		return
	}

	_, _ = description.Printf("Read File List from %s\n", path)
	start := time.Now()

	wl := make(filter.File)
	bl := make(filter.File)
	if confFilter != nil {
		wl, _ = confFilter.File(filter.KeyWhitelist)
		bl, _ = confFilter.File(filter.KeyBlacklist)
	}

	fileList, err := mp3files.GetFileList(path, wl, bl)
	if err != nil {
		res <- FileChannel{Err: err}
		return
	}

	duration(start, time.Now(), fmt.Sprintf("%d files found in %s", len(fileList), path))

	res <- FileChannel{Files: fileList}
}
