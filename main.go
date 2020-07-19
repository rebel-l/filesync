// TODO
package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"

	"github.com/fatih/color"

	"github.com/rebel-l/go-utils/osutils"
	"github.com/rebel-l/go-utils/pb"
	"github.com/rebel-l/mp3sync/mp3files"
	"github.com/rebel-l/mp3sync/sync"
	"github.com/rebel-l/mp3sync/transform"
)

const (
	source      = "D:\\CarMusic - Prepare" // TODO: config file
	destination = "F:\\"                   // TODO: config file
)

var (
	errPathNotExisting = errors.New("path does not exist")
	errFormat          = color.New(color.FgRed)
	description        = color.New(color.FgGreen)  // nolint: gochecknoglobals
	listFormat         = color.New(color.FgHiBlue) // nolint: gochecknoglobals
)

func main() {
	title := color.New(color.Bold, color.FgGreen)
	_, _ = title.Println("MP3 sync started ...")
	fmt.Println()

	info := color.New(color.FgYellow)

	fmt.Printf("%s %s\n", description.Sprint("Source:"), info.Sprint(source))
	fmt.Printf("%s %s\n", description.Sprint("Destination:"), info.Sprint(destination))

	fmt.Println()

	if err := do(); err != nil {
		fmt.Println()

		_, _ = errFormat.Printf("MP3 sync finished with error: %v\n", err)
	} else {
		fmt.Println()

		_, _ = title.Println("MP3 sync finished successful!")
	}
}

func do() error {
	if !osutils.FileOrPathExists(source) {
		return fmt.Errorf("%w: %s", errPathNotExisting, source)
	}

	if !osutils.FileOrPathExists(destination) {
		return fmt.Errorf("%w: %s", errPathNotExisting, destination)
	}

	fileList, err := readFileList()
	if err != nil {
		return err
	}

	syncFiles, err := diff(fileList)
	if err != nil {
		return err
	}

	listDiff(syncFiles)

	return snycFiles(syncFiles)
}

func readFileList() (mp3files.Files, error) {
	_, _ = description.Print("Read File List: ")
	start := time.Now()

	defer fmt.Println()

	fileList, err := mp3files.GetFileList(source)
	if err != nil {
		return nil, err
	}

	duration(start, time.Now(), fmt.Sprintf("%d files found", len(fileList)))

	return fileList, nil
}

func duration(start, finish time.Time, msg string) {
	_, _ = description.Printf("%s in %s\n", msg, finish.Sub(start))
}

func diff(fileList mp3files.Files) (sync.Files, error) {
	_, _ = description.Println("Analyse files to be synced ...")
	start := time.Now()

	defer fmt.Println()

	bar := pb.New(pb.EngineCheggaaa, len(fileList))

	var syncFiles sync.Files

	for _, v := range fileList {
		bar.Increment()

		destinatonFileName, err := transform.Do(destination, v)
		if err != nil {
			// return nil, err // TODO: add errors to stack and log to file at the end
			continue
		}

		f := sync.File{Source: v, Destination: destinatonFileName}

		inSync, err := f.IsInSync()
		if err != nil {
			return nil, err
		}

		if !inSync {
			syncFiles = append(syncFiles, f)
		}
	}

	bar.Finish()

	duration(start, time.Now(), fmt.Sprintf("%d files analysed", len(fileList)))

	return syncFiles, nil
}

func listDiff(files sync.Files) {
	_, _ = listFormat.Printf("Total files to sync: %d\n", len(files))

	t := prompt.Input("Show Diff? [Y/n] ", func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix([]prompt.Suggest{}, d.GetWordBeforeCursor(), true)
	})

	if strings.ToLower(t) != "n" {
		for _, v := range files {
			_, _ = listFormat.Println(v.Destination)
		}
	}

	fmt.Println()
}

func snycFiles(files sync.Files) error {
	t := prompt.Input("Start Sync? [Y/n] ", func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix([]prompt.Suggest{}, d.GetWordBeforeCursor(), true)
	})

	if strings.ToLower(t) == "n" {
		return errors.New("aborted by user")
	}

	_, _ = description.Print("Sync files: ")
	start := time.Now()

	defer fmt.Println()

	bar := pb.New(pb.EngineCheggaaa, len(files))

	for _, v := range files {
		bar.Increment()

		if err := sync.Do(v); err != nil {
			return err // TODO: add errors to stack and log to file at the end
		}
	}

	bar.Finish()

	duration(start, time.Now(), fmt.Sprintf("%d files synced", len(files)))

	return nil
}

// TODO:
// 1. Check Space needed / left
// 2. Collect errors / Write to log or display
// 3. Use Config: source, destination, filters
// 4. Documentation
// 5. Tests
