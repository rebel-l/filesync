package config

import "github.com/rebel-l/mp3sync/filter"

type Config struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	// TODO: Whitelist / Blacklist can be managed as map
	WhiteList filter.Filter `json:"whiteList"`
	BlackList filter.Filter `json:"blackList"`
}
