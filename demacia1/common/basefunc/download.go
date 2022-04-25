package basefunc

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

func Download(i string) string {
	u, err := url.Parse(i)
	if err != nil {
		return ""
	}
	ext := path.Ext(u.Path)
	if ext == "" || ext == "." {
		return ""
	}
	p := os.TempDir() + "/" + RandNumberString() + ext
	res, err := http.Get(i)
	if err != nil {
		return ""
	}
	file, err := os.Create(p)
	if err != nil {
		return ""
	}
	_, _ = io.Copy(file, res.Body)
	return p
}
