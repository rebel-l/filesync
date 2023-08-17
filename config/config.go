package config

type Config struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	// TODO: Whitelist / Blacklist can be managed as map
	WhiteList Filter `json:"whiteList"`
	BlackList Filter `json:"blackList"`
}
