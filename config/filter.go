package config

type Filter struct {
	File `json:"file"`
}

//func (f Filter) Contains(content string) bool {
//	for _, v := range f {
//		if strings.Contains(content, v) {
//			return true
//		}
//	}
//
//	return false
//}
