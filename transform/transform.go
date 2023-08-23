package transform

import (
	"errors"
	"fmt"
	"github.com/bogem/id3v2/v2"
	"github.com/rebel-l/mp3sync/filter"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rebel-l/go-utils/pb"
	"github.com/rebel-l/mp3sync/mp3files"
)

const (
	defaultSubfolder = "default"
	numericSubfolder = "#"
	frameIDDisk      = "TPOS"
	frameIDTrack     = "TRCK"
	MP3Extension     = ".mp3"
)

var ErrParseTag = errors.New("failed to parse mp3 tag")

func Do(fileList mp3files.Files, destination string, source string, whiteList filter.MP3Tag, blackList filter.MP3Tag) (mp3files.Files, []error) {
	var errs []error
	var newFiles mp3files.Files

	bar := pb.New(pb.EngineCheggaaa, len(fileList))
	defer bar.Finish()

	for _, v := range fileList {
		bar.Increment()

		if filepath.Ext(v.Info.Name()) == MP3Extension {
			var err error
			v.Name, err = transform(destination, source, v, whiteList, blackList)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			if v.Name != "" {
				newFiles = append(newFiles, v)
			}
		}
	}

	return newFiles, errs
}

func transform(destination string, source string, f mp3files.File, whiteList filter.MP3Tag, blackList filter.MP3Tag) (string, error) {
	tag, err := loadTag(f) // TODO: should be outside of this package as it is not part of transformer
	if err != nil {
		return "", fmt.Errorf("%w from %s: %v", ErrParseTag, f.Name, err)
	}

	if (len(whiteList) > 0 && !whiteList.Contains(tag)) || (len(blackList) > 0 && blackList.Contains(tag)) {
		return "", nil
	}

	name, err := getFileName(tag, filepath.Ext(f.Info.Name()))
	if err != nil {
		return "", fmt.Errorf("%w from %s: %v", ErrParseTag, f.Name, err)
	}

	return filepath.Join(destination, getSubFolder(f.Name, source), name), nil
}

func getSubFolder(fileName string, source string) string {
	subFolder := defaultSubfolder

	source = strings.Replace(fileName, source+string(os.PathSeparator), "", 1)

	parts := strings.Split(source, string(os.PathSeparator))

	if len(parts) > 0 {
		subFolder = strings.ToUpper(string(parts[0][0]))
	}

	match, _ := regexp.MatchString("[A-Z]", subFolder)
	if !match {
		subFolder = numericSubfolder
	}

	return subFolder
}

func getFileName(tag *id3v2.Tag, ext string) (string, error) {

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

	return replaceChars(name + " - " + tag.Title() + ext), nil
}

func replaceChars(s string) string {
	chars := map[string]string{
		":":  ";",
		"\\": "",
		"/":  "",
		"?":  "¿",
		"\"": ",",
		"'":  ",",
		"*":  "x",
		"+":  "x",
		"[":  "(",
		"]":  ")",
		">":  "-",
		"<":  "-",
		"|":  "-",
	}

	for k, v := range chars {
		s = strings.Replace(s, k, v, -1)
	}

	return s
}
