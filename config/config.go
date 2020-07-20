package config

type Config struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Filter      Filter `json:"filter"`
}
