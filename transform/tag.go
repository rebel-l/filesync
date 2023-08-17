package transform

import (
	"github.com/bogem/id3v2/v2"
	"github.com/rebel-l/mp3sync/mp3files"
)

func loadTag(f mp3files.File) (*id3v2.Tag, error) {
	tag, err := id3v2.Open(f.Name, id3v2.Options{Parse: true})
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tag.Close()
	}()

	return tag, nil
}
