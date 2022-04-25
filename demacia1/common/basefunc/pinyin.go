package basefunc

import (
	"github.com/mozillazg/go-pinyin"
)

func GetPinyin(name string) string {
	return pinyin.Slug(name, pinyin.Args{
		Fallback: func(r rune, a pinyin.Args) []string {
			return []string{string(r)}
		}})
}
