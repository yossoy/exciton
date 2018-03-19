// +build !release

package markup

import (
	"path/filepath"
)

func (k *Klass) getResourcePath(base string) string {
	return filepath.Join(k.dir, "resources")
}
