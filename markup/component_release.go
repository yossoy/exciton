// +build release

package markup

import (
	"net/http"
	"path"
	"path/filepath"

	"github.com/yossoy/exciton/driver"
)

func (k *Klass) GetResourceFile(fn string) (http.File, error) {
	fs, err := driver.ResourcesFileSystem()
	if err != nil {
		return nil, err
	}
	fp := path.Join(k.pathInfo.pkgPath, fn)
	return fs.Open(fp)
}

func (k *Klass) getResourcePath(base string) string {
	return filepath.Join(base, filepath.FromSlash(k.pathInfo.pkgPath))
}

// func GetComponentCSSFiles(resPath string) ([]string, error) {
// 	var cssFiles []string
// 	for _, k := range componentKlasses {
// 		if k.cssFile != "" {
// 			basePath := k.getResourcePath(resPath)
// 			cssPath := filepath.Join(basePath, k.cssFile)
// 			relPath, err := filepath.Rel(resPath, cssPath)
// 			if err != nil {
// 				return nil, err
// 			}
// 			cssFiles = append(cssFiles, relPath)
// 		}
// 	}
// 	return cssFiles, nil
// }

// func GetComponentJSFiles(resPath string) ([]string, error) {
// 	var jsFiles []string
// 	for _, k := range componentKlasses {
// 		if k.jsFile != "" {
// 			basePath := k.getResourcePath(resPath)
// 			jsPath := filepath.Join(basePath, k.jsFile)
// 			relPath, err := filepath.Rel(resPath, jsPath)
// 			if err != nil {
// 				return nil, err
// 			}
// 			jsFiles = append(jsFiles, relPath)
// 		}
// 	}
// 	return jsFiles, nil
// }
