package filter

import (
	"github.com/bogem/id3v2/v2"
	"strings"
)

const (
	KeyAlbum  = "album"
	KeyArtist = "artist"
	KeyGenre  = "genre"
)

var (
	possibleTagKeys = []Key{KeyAlbum, KeyArtist, KeyGenre}
)

type MP3Tag map[Key][]string

func (t MP3Tag) Contains(tag *id3v2.Tag) bool {
	for _, key := range possibleTagKeys {
		patterns, ok := t[key]
		if !ok {
			continue
		}

		for _, pattern := range patterns {
			if matchTag(key, tag, pattern) {
				return true
			}
		}
	}

	return false
}

func matchTag(key Key, tag *id3v2.Tag, pattern string) bool {
	switch key {
	case KeyAlbum:
		return strings.Contains(tag.Album(), pattern)
	case KeyArtist:
		return strings.Contains(tag.Artist(), pattern)
	case KeyGenre:
		return strings.Contains(tag.Genre(), pattern)
	}

	return false
}
