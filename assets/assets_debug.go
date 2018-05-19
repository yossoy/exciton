// +build !release

package assets

import (
	"net/http"
	"path/filepath"
	"runtime"
)

var FileSystem = func() http.FileSystem {
	_, fp, _, _ := runtime.Caller(0)
	return http.Dir(filepath.Join(filepath.Dir(fp), "data"))
}()
