package filesync

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rebel-l/go-utils/osutils"
	"github.com/rebel-l/go-utils/pb"
)

var ErrCreateDirectory = errors.New("failed to create destination directory")

func Do(files Files) []error {
	var errs []error

	defer fmt.Println()

	bar := pb.New(pb.EngineCheggaaa, len(files))

	for _, v := range files {
		bar.Increment()

		var operation func(file File) error
		switch v.Operation {
		case OperationCreate, OperationUpdate:
			operation = copy
			break
		case OperationDelete:
			operation = delete
			break
		default:
			errs = append(errs, fmt.Errorf("operation %q for file %q / %q not known", v.Operation, v.Source.Name, v.Destination.Name))
			continue
		}

		if err := operation(v); err != nil {
			errs = append(errs, err)
		}
	}

	bar.Finish()

	return errs
}

func copy(file File) error {
	destPath, _ := filepath.Split(file.Destination.Name)
	if err := osutils.CreateDirectoryIfNotExists(destPath); err != nil {
		return fmt.Errorf("%w - %q: %v", ErrCreateDirectory, destPath, err)
	}

	if err := osutils.CopyFile(file.Source.Name, file.Destination.Name); err != nil {
		return err
	}

	return os.Chtimes(file.Destination.Name, file.Source.Info.ModTime(), file.Source.Info.ModTime())
}

func delete(file File) error {
	//return os.Remove(file.Destination.Name)
	return fmt.Errorf("DELETE: %s", file.Destination.Name) // TODO: remove this line and activate the other one
}
