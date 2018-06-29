// +build !release

package windows

/*
#cgo CFLAGS: -DDEBUG -g
#cgo CXXFLAGS: -DDEBUG -g
*/
import "C"

import (
	"net/http"
	"os"
	"path/filepath"
)

func (d *windows) ResourcesFileSystem() (http.FileSystem, error) {
	exePathStr, err := os.Executable()
	if err != nil {
		return nil, err
	}
	resourcesName := filepath.Join(filepath.Dir(exePathStr), "resources")
	return http.Dir(resourcesName), nil
}
