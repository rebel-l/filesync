package filter

const (
	KeyBlacklist = "blackList"
	KeyWhitelist = "whiteList"
)

type BlackWhiteList map[string]Filter

func (l BlackWhiteList) Blacklist() (Filter, bool) {
	f, ok := l[KeyBlacklist]
	return f, ok
}

func (l BlackWhiteList) Whitelist() (Filter, bool) {
	f, ok := l[KeyWhitelist]
	return f, ok
}

func (l BlackWhiteList) File(blackOrWhitelist string) (File, bool) {
	f, ok := l[blackOrWhitelist]
	if !ok {
		return nil, false
	}

	return f.File, true
}

func (l BlackWhiteList) MP3Tag(blackOrWhitelist string) (MP3Tag, bool) {
	f, ok := l[blackOrWhitelist]
	if !ok {
		return nil, false
	}

	return f.MP3Tag, true
}
