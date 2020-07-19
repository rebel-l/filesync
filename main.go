// TODO
package main

import (
	"errors"
	"flag"
	"fmt"

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
	previewFlag        *bool // nolint: gochecknoglobals
)

func main() {
	flags()

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
	if err != nil {
		return err
	}

	var bar pb.Progressor
	if *previewFlag {
		bar = pb.New(pb.EngineBlackhole, len(fileList))
	} else {
		bar = pb.New(pb.EngineCheggaaa, len(fileList))
	}

	defer bar.Finish()

	for _, v := range fileList {
		bar.Increment()

		destinatonFileName, err := transform.Do(destination, v)
		if err != nil {
			return err // TODO: add errors to stack and log to file at the end
		}

		if err := sync.Do(v, destinatonFileName, *previewFlag); err != nil {
			return err // TODO: add errors to stack and log to file at the end
		}
	}

	return nil
}

func flags() {
	previewFlag = flag.Bool("p", false, "shows a preview of the files to synced to destination path")
	flag.Parse()
}
