// TODO
package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/rebel-l/go-utils/osutils"
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

	errFormat = color.New(color.FgRed)
)

func main() {
	title := color.New(color.Bold, color.FgGreen)
	_, _ = title.Println("MP3 sync started ...")
	fmt.Println()

	description := color.New(color.FgGreen)
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

	fileList, err := mp3files.GetFileList(source)
	if err != nil && !strings.Contains(err.Error(), "100 files reached") {
		// TODO: remove temporary check for file limit check
		return err
	}

	// TODO: progress bar
	for _, v := range fileList {
		destinatonFileName, err := transform.Do(destination, v)
		if err != nil {
			return err // TODO: add errors to stack and log to file at the end
		}

		if err := sync.Do(v, destinatonFileName, true); err != nil { // TODO: preview should be controlled by script parameter
			return err // TODO: add errors to stack and log to file at the end
		}
	}

	return nil
}
