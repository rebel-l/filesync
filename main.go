// TODO
package main

import (
	"errors"
	"fmt"
	"strings"

	"gihub.com/rebel-l/mp3sync/mp3files"

	"github.com/fatih/color"

	"github.com/rebel-l/go-utils/osutils"
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
		_, _ = errFormat.Printf("MP3 sync finished with error: %v\n", err)
	} else {
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

	listFormat := color.New(color.FgHiBlue)
	// TODO: progress bar
	for _, v := range fileList {
		_, _ = listFormat.Println(v.GetName())
	}

	return nil
}
