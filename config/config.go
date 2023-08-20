package config

import "github.com/rebel-l/mp3sync/filter"

type Config struct {
	Source      string                `json:"source"`
	Destination string                `json:"destination"`
	Filter      filter.BlackWhiteList `json:"filter"`
}
