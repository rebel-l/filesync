package filter

const (
	KeyBlacklist = "blacklist"
	KeyWhitelist = "whitelist"
)

type BlackWhiteList map[string]Filter

func (l BlackWhiteList) Blacklist() (Filter, bool) { // TODO: use key to get File / MP3Tag
	f, ok := l[KeyBlacklist]
	return f, ok
}

func (l BlackWhiteList) Whitelist() (Filter, bool) { // TODO: use key to get File / MP3Tag
	f, ok := l[KeyWhitelist]
	return f, ok
}
