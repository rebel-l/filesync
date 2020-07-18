package transform

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bogem/id3v2"

	"github.com/rebel-l/mp3sync/mp3files"
)

const (
	defaultSubfolder  = "default"
	numericSubfolder  = "#"
	minDirectoryParts = 2
	frameIDDisk       = "TPOS"
	frameIDTrack      = "TRCK"
)

var (
	ErrParseTag = errors.New("failed to parse mp3 tag")
)

func Do(path string, f mp3files.File) (string, error) {
	name, err := getFileName(f)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrParseTag, err)
	}

	return filepath.Join(path, getSubFolder(f.Path), name), nil
}

func getSubFolder(path string) string {
	subFolder := defaultSubfolder

	parts := strings.Split(path, string(os.PathSeparator))
	if len(parts) > minDirectoryParts {
		subFolder = strings.ToUpper(string(parts[2][0]))
	}

	match, _ := regexp.MatchString("[A-Z]", subFolder)
	if !match {
		subFolder = numericSubfolder
	}

	return subFolder
}

func getFileName(f mp3files.File) (string, error) {
	tag, err := id3v2.Open(f.GetName(), id3v2.Options{Parse: true})
	if err != nil {
		return "", err
	}

	defer func() {
		_ = tag.Close()
	}()

	name := tag.Artist()

	if tag.Album() != "" {
		name += " - " + tag.Album()
	}

	if tag.Year() != "" {
		name += " (" + tag.Year() + ")"
	}

	disk := tag.GetTextFrame(frameIDDisk).Text
	if disk != "" {
		name += " - " + disk
	}

	track := tag.GetTextFrame(frameIDTrack).Text
	if track != "" {
		if len(track) == 1 {
			track = "0" + track
		}

		name += " - " + track
	}

	return replaceChars(name + " - " + tag.Title() + mp3files.Extension), nil
}

func replaceChars(s string) string {
	chars := map[string]string{
		":":  ";",
		"\\": "",
		"/":  "",
		"?":  "",
	}

	for k, v := range chars {
		s = strings.Replace(s, k, v, -1)
	}

	return s
}
