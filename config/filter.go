package config

import "strings"

type Filter []string

func (f Filter) Contains(content string) bool {
	for _, v := range f {
		if strings.Contains(content, v) {
			return true
		}
	}

	return false
}
