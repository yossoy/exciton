// +build !release

package markup

import (
	"net/http"
	"path"
	"path/filepath"
)

func (k *klassPathInfo) getResourceFile(fn string) (http.File, error) {
	fp := path.Join("resources", fn)
	return http.Dir(k.dir).Open(fp)
}

func (k *klassPathInfo) getResourcePath(base string) string {
	return filepath.Join(k.dir, "resources")
}
