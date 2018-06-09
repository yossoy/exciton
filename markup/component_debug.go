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

// var (
// 	basePathToIDMap = make(map[string]string)
// 	idToBasePathMap = make(map[string]string)
// )

// func pathToComponentBaseID(k *Klass) string {
// 	p := k.Path
// 	if id, ok := basePathToIDMap[p]; ok {
// 		return id
// 	}
// 	id := fmt.Sprintf("id%d", len(basePathToIDMap))
// 	basePathToIDMap[p] = id
// 	idToBasePathMap[id] = p
// 	return id
// }

// func componentBaseIDToPath(id string) (string, bool) {
// 	path, ok := idToBasePathMap[id]
// 	return path, ok
// }

// func getComponentCSSJSFile(k *Klass, resPath string, file string, goext string) ([]byte, error) {
// 	basePath := k.getResourcePath(resPath)
// 	cssPath := filepath.Join(basePath, k.cssFile)
// 	relPath, err := filepath.Rel(resPath, basePath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var b []byte
// 	if filepath.Ext(k.cssFile) == goext {
// 		b, err = ReadComponentNamespaceFile(filepath.ToSlash(relPath), cssPath, k.Path)
// 	} else {
// 		b, err = ioutil.ReadFile(cssPath)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	return b, nil
// }

// func GetComponentCSSFiles(resPath string) ([]string, error) {
// 	cssFiles := make([]string, 0, len(componentKlasses))
// 	for _, k := range componentKlasses {
// 		if k.cssFile != "" {
// 			b, err := getComponentCSSJSFile(k, resPath, k.cssFile, ".gocss")
// 			if err != nil {
// 				return nil, err
// 			}
// 			cssFiles = append(cssFiles, string(b))
// 		}
// 	}
// 	return cssFiles, nil
// }

// func GetComponentJSFiles(resPath string) ([]string, error) {
// 	var jsFiles []string
// 	for _, k := range componentKlasses {
// 		if k.jsFile != "" {
// 			b, err := getComponentCSSJSFile(k, resPath, k.jsFile, ".gojs")
// 			if err != nil {
// 				return nil, err
// 			}
// 			jsFiles = append(jsFiles, string(b))
// 		}
// 	}
// 	return jsFiles, nil
// }

// func GetComponentWrappedCSS(css string) ([]byte, error) {
// 	var bb bytes.Buffer
// 	scn := scanner.New(css)
// }
