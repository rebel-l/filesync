package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rebel-l/go-utils/osutils"
)

const (
	errFileNotFound = "configuration file not found"
	errLoadFile     = "failed to load configuration"
	errNoJSONFormat = "content of file is not in JSON format"
	errReadData     = "failed to read data from file"
)

var (
	// ErrFileNotFound is the error if the config file doesn't exist
	ErrFileNotFound = errors.New(errFileNotFound)

	// ErrNoJSONFormat indicates that the content is not a JSON
	ErrNoJSONFormat = errors.New(errNoJSONFormat)

	ErrLoadFile = errors.New(errLoadFile)

	ErrReadData = errors.New(errReadData)
)

func Load(fileName string) (*Config, error) {
	fileName = filepath.Clean(fileName)
	if !osutils.FileOrPathExists(fileName) {
		return nil, ErrFileNotFound
	}

	f, err := os.Open(fileName) // nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLoadFile, err)
	}

	defer func() {
		_ = f.Close()
	}()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadData, err)
	}

	c := &Config{}
	if err = json.Unmarshal(data, c); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNoJSONFormat, err)
	}

	return c, nil
}
