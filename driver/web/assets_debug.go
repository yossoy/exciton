// +build !release

package web

import (
	"net/http"
	"path/filepath"
	"runtime"
)

var fileSystem = func() http.FileSystem {
	_, fp, _, _ := runtime.Caller(0)
	return http.Dir(filepath.Join(filepath.Dir(fp), "data"))
}()
