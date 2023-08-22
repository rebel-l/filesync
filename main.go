package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/rebel-l/go-utils/osutils"
	"github.com/rebel-l/go-utils/pb"
	"github.com/rebel-l/mp3sync/config"
	"github.com/rebel-l/mp3sync/filter"
	"github.com/rebel-l/mp3sync/mp3files"
	"github.com/rebel-l/mp3sync/sync"
	"github.com/rebel-l/mp3sync/transform"
	"github.com/shirou/gopsutil/v3/disk"
)

const (
	configFile        = "config.json"
	logPath           = "logs"
	logFileNameFormat = "20060102-150405"
	MP3Extension      = ".mp3"
)

var (
	errPathNotExisting = errors.New("path does not exist")
	errWriteLog        = errors.New("failed to write log file")
	errAbortedByUser   = errors.New("aborted by user")
	errFormat          = color.New(color.FgRed)
	description        = color.New(color.FgGreen)  // nolint: gochecknoglobals
	listFormat         = color.New(color.FgHiBlue) // nolint: gochecknoglobals
)

func main() {
	title := color.New(color.Bold, color.FgGreen)
	_, _ = title.Println("MP3 sync started ...")
	fmt.Println()

	info := color.New(color.FgYellow)

	conf, err := config.Load(filepath.Join(".", configFile))
	if err != nil {
		_, _ = errFormat.Printf("failed to load config: %v", err)
		return
	}

	_, _ = description.Printf("Source: %s\n", info.Sprint(conf.Source))
	_, _ = description.Printf("Destination: %s\n", info.Sprint(conf.Destination))

	fmt.Println()

	if err := do(conf); err != nil {
		fmt.Println()

		_, _ = errFormat.Printf("MP3 sync finished with error: %v\n", err)
	} else {
		fmt.Println()

		_, _ = title.Println("MP3 sync finished successful!")
	}
}

func do(conf *config.Config) error {
	// TODO:
	// 1. move to package and call it in 2 go routines with source (incl. filter) and destination (excl. filter)
	// 2. transform source
	// 3. diff file sizes + source / destination and set operations: create / update / delete
	// 4. ask to list diff?
	// 5. run operations
	sourceList, err := readFileList(conf.Source, conf.Filter)
	if err != nil {
		return err
	}
	fmt.Println("SOURCE", sourceList)

	destinationList, err := readFileList(conf.Destination, nil)
	if err != nil {
		return err
	}
	fmt.Println("DESTINATION", destinationList)

	var globErr bool

	syncFiles, errs := diff(sourceList, conf.Destination, conf.Source, conf.Filter)
	if len(errs) > 0 {
		globErr = true

		if err := showAndLogErrors(errs); err != nil {
			return err
		}
	}

	listDiff(syncFiles)

	if err := diskSpace(syncFiles, conf.Destination); err != nil {
		return err
	}

	errs = snycFiles(syncFiles)
	if len(errs) > 1 {
		if errors.Is(errs[0], errAbortedByUser) {
			return errs[0]
		}

		globErr = true

		logFileName, err := logErrors(errs)
		if err != nil {
			return err
		}

		_, _ = errFormat.Printf("found %d errors\n", len(errs))
		_, _ = errFormat.Printf("logged errors in file %s\n", logFileName)
	}

	if globErr {
		return errors.New("see log for more details")
	}

	return nil
}

func duration(start, finish time.Time, msg string) {
	_, _ = description.Printf("%s in %s\n", msg, finish.Sub(start))
}

func diff(fileList mp3files.Files, destination string, source string, confFilter filter.BlackWhiteList) (sync.Files, []error) {
	_, _ = description.Println("Analyse files to be synced ...")
	start := time.Now()

	defer fmt.Println()

	bar := pb.New(pb.EngineCheggaaa, len(fileList))

	var syncFiles sync.Files

	var errs []error

	wl, _ := confFilter.MP3Tag(filter.KeyWhitelist)
	bl, _ := confFilter.MP3Tag(filter.KeyBlacklist)

	for _, v := range fileList {
		bar.Increment()

		destinationFileName := filepath.Join(destination, v.Name) // FIXME v.Name includes full path and doesn't work like this
		if filepath.Ext(v.Info.Name()) == MP3Extension {
			var err error
			destinationFileName, err = transform.Do(destination, source, v, wl, bl)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			if destinationFileName == "" {
				continue
			}
		}

		destFile := mp3files.File{Name: destinationFileName}

		if osutils.FileOrPathExists(destinationFileName) {
			var err error
			destFile.Info, err = os.Lstat(destinationFileName)

			if err != nil {
				errs = append(errs, err)
				continue
			}
		}

		f := sync.File{Source: v, Destination: destFile}

		if !f.IsInSync() {
			syncFiles = append(syncFiles, f)
		}
	}

	bar.Finish()

	duration(start, time.Now(), fmt.Sprintf("%d files analysed", len(fileList)))

	return syncFiles, errs
}

func listDiff(files sync.Files) {
	_, _ = listFormat.Printf("Total files to sync: %d\n", len(files))

	t := prompt.Input("Show Diff? [Y/n] ", func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix([]prompt.Suggest{}, d.GetWordBeforeCursor(), true)
	})

	if strings.ToLower(t) != "n" {
		for _, v := range files {
			_, _ = listFormat.Println(v.Destination.Name)
		}
	}

	fmt.Println()
}

func snycFiles(files sync.Files) []error {
	var errs []error

	t := prompt.Input("Start Sync? [Y/n] ", func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix([]prompt.Suggest{}, d.GetWordBeforeCursor(), true)
	})

	if strings.ToLower(t) == "n" {
		errs = append(errs, errAbortedByUser)
		return errs
	}

	_, _ = description.Print("Sync files: ")
	start := time.Now()

	defer fmt.Println()

	bar := pb.New(pb.EngineCheggaaa, len(files))

	for _, v := range files {
		bar.Increment()

		if err := sync.Do(v); err != nil {
			errs = append(errs, err)
		}
	}

	bar.Finish()

	duration(start, time.Now(), fmt.Sprintf("%d files synced", len(files)))

	return errs
}

func diskSpace(files sync.Files, destination string) error {
	di, err := disk.Usage(destination)
	if err != nil {
		return err
	}

	needed := files.SpaceNeeded()
	left := int64(di.Free) - needed

	neededDisplay := humanize.Bytes(uint64(needed))
	if needed < 0 {
		neededDisplay = "-" + humanize.Bytes(uint64(needed*-1))
	}

	leftDisplay := humanize.Bytes(uint64(left))
	if left < 0 {
		leftDisplay = "-" + humanize.Bytes(uint64(left*-1))
	}

	_, _ = listFormat.Printf("Free Disk Space: %s\n", humanize.Bytes(di.Free))
	_, _ = listFormat.Printf("Disk Space Needed: %s\n", neededDisplay)
	_, _ = listFormat.Printf("Disk Space Left: %s\n", leftDisplay)

	fmt.Println()

	if left < 1 {
		return fmt.Errorf("not enough free disk space, need %d more bytes", left*-1)
	}

	return nil
}

func showAndLogErrors(errs []error) error {
	logFileName, err := logErrors(errs)
	if err != nil {
		return err
	}

	return showErrors(logFileName, errs)
}

func showErrors(logFileName string, errs []error) error {
	_, _ = errFormat.Printf("found %d errors\n", len(errs))

	t := prompt.Input("Continue (errored files will be skipped)? [Y/n/s = show files] ",
		func(d prompt.Document) []prompt.Suggest {
			return prompt.FilterHasPrefix([]prompt.Suggest{}, d.GetWordBeforeCursor(), true)
		},
	)

	_, _ = errFormat.Printf("logged errors in file %s\n", logFileName)

	switch strings.ToLower(t) {
	case "n":
		return errAbortedByUser
	case "s":
		for _, e := range errs {
			_, _ = errFormat.Println(e.Error())
		}
	}

	return nil
}

func logErrors(errs []error) (string, error) {
	logFileName := filepath.Join(".", logPath, fmt.Sprintf("%s.log", time.Now().Format(logFileNameFormat)))

	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND, 0o600)
	if err != nil {
		return logFileName, fmt.Errorf("%w: %v", errWriteLog, err)
	}

	defer func() {
		_ = logFile.Close()
	}()

	for _, e := range errs {
		if _, err := logFile.WriteString(e.Error() + "\n"); err != nil {
			return logFileName, fmt.Errorf("%w: %v", errWriteLog, err)
		}
	}

	return logFileName, nil
}

// TODO:
// 2. delete not matching files
// 3. activate all linters
// 4. Documentation / Badges: licence, goreport, issues, releases
// 5. Tests / Badges: build, coverage
