package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/rebel-l/go-utils/pb"
	"github.com/rebel-l/mp3sync/filesync"
	"github.com/rebel-l/mp3sync/mp3files"
)

func diff(source, destination mp3files.Files) filesync.Files {
	_, _ = description.Println("Create diff of files to be synced ...")
	start := time.Now()

	defer fmt.Println()

	var syncFiles filesync.Files

	count := len(source) + len(destination)
	bar := pb.New(pb.EngineCheggaaa, count)
	defer bar.Finish()

	for _, s := range source {
		bar.Increment()
		d, ok := destination[s.Name]

		// if destination doesn't exist just add it
		if !ok {
			syncFiles = append(syncFiles, filesync.File{Source: s.Info, Destination: filesync.FileInfo{FileName: s.Name}, Operation: filesync.OperationCopy})
			continue
		}

		// if destination exists it must be out of sync
		newFile := filesync.File{Source: s.Info, Destination: d.Info, Operation: filesync.OperationCopy}
		if !newFile.IsInSync() {
			syncFiles = append(syncFiles, newFile)
		}
	}

	for _, d := range destination {
		bar.Increment()
		// files not in source needs to be deleted
		_, ok := source[d.Name]
		if !ok {
			syncFiles = append(syncFiles, filesync.File{Destination: d.Info, Operation: filesync.OperationDelete})
		}
	}

	duration(start, time.Now(), fmt.Sprintf("diff for %d files created", count))

	return syncFiles
}

func listDiff(files filesync.Files) {
	_, _ = listFormat.Printf("Total files to sync: %d\n", len(files))

	t := prompt.Input("Show Diff? [Y/n] ", func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix([]prompt.Suggest{}, d.GetWordBeforeCursor(), true)
	})

	if strings.ToLower(t) != "n" {
		for _, v := range files {
			_, _ = listFormat.Printf("%s: %s\n", v.Operation, v.Destination.Name())
		}
	}

	fmt.Println()
}
