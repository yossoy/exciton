// +build !release

package markup

import (
	"net/http"
	"path"
	"path/filepath"
)

func (k *Klass) GetResourceFile(fn string) (http.File, error) {
	fp := path.Join("resources", fn)
	return http.Dir(k.pathInfo.dir).Open(fp)
}

func (k *Klass) getResourcePath(base string) string {
	return filepath.Join(k.pathInfo.dir, "resources")
}
