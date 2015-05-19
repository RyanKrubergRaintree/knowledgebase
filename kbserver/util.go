package kbserver

import (
	"path"
	"strings"
)

func SafeFile(dir string, url string) string {
	upath := url
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
	}
	upath = path.Clean(upath)
	return path.Join(dir, upath[1:])
}
