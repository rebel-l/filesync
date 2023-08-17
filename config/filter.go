package config

// TODO: move to own package filter

type Filter struct {
	File   `json:"file"`
	MP3Tag Tag `json:"mp3tag"`
}
