// +build release

package markup

import (
	"net/http"
	"path"
	"path/filepath"

	"github.com/yossoy/exciton/driver"
)

func (k *klassPathInfo) getResourceFile(fn string) (http.File, error) {
	fs, err := driver.ResourcesFileSystem()
	if err != nil {
		return nil, err
	}
	fp := path.Join(k.pkgPath, fn)
	return fs.Open(fp)
}

func (k *klassPathInfo) getResourcePath(base string) string {
	return filepath.Join(base, filepath.FromSlash(k.pkgPath))
}
