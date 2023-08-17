package config

import (
	"os"
	"path/filepath"
	"strings"
)

// TODO: move to own package filter

const (
	KeyName      = "name"
	KeyExtension = "ext"
)

var (
	possibleKeys = []Key{KeyName, KeyExtension}
)

type Key string

type File map[Key][]string

func (f File) Contains(info os.FileInfo) bool {
	for _, key := range possibleKeys {
		patterns, ok := f[key]
		if !ok {
			continue
		}

		for _, pattern := range patterns {
			if match(key, info, pattern) {
				return true
			}
		}
	}

	return false
}

func match(key Key, info os.FileInfo, pattern string) bool {
	switch key {
	case KeyName:
		return strings.Contains(info.Name(), pattern)
	case KeyExtension:
		return strings.ToLower(filepath.Ext(info.Name())) == pattern
	}

	return false
}
