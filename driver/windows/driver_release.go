// +build release

package windows

/*
#cgo CFLAGS: -DNDEBUG
#cgo CXXFLAGS: -DNDEBUG

#include "driver.h"
*/
import "C"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/yossoy/exciton/driver/windows/resfile"
)

var resRootItem *resfile.Item

func initResourceFileSystem() error {
	b, err := getResourceFile(resfile.FileMapJsonID)
	if err != nil {
		panic("initResourceFile: getResource failed:")
		return err
	}
	var resItems resfile.Items
	err = json.Unmarshal(b, &resItems)
	if err != nil {
		return err
	}
	resRootItem = &resfile.Item{
		ID:       0,
		Name:     "",
		Size:     0,
		Mode:     os.ModeDir,
		ModTime:  time.Now(),
		IsDir:    true,
		Children: resItems,
	}
	return nil
}

func getResourceFile(no int) ([]byte, error) {
	rc := C.Driver_GetResFile(C.int(no))
	if rc.ptr == nil {
		return nil, fmt.Errorf("Resource[%d] not found", no)
	}
	return C.GoBytes(rc.ptr, rc.size), nil
}

type resFileStat struct{ item *resfile.Item }

func (fs *resFileStat) Name() string       { return fs.item.Name }
func (fs *resFileStat) Size() int64        { return fs.item.Size }
func (fs *resFileStat) Mode() os.FileMode  { return fs.item.Mode }
func (fs *resFileStat) ModTime() time.Time { return fs.item.ModTime }
func (fs *resFileStat) IsDir() bool        { return fs.item.IsDir }
func (fs *resFileStat) Sys() interface{}   { return fs.item }

type resFile struct {
	item   *resfile.Item
	reader *bytes.Reader
}

func (f *resFile) Close() error {
	if f.reader != nil {
		f.reader = nil
	}
	return nil
}
func (f *resFile) Read(b []byte) (int, error) {
	if f.reader != nil {
		return f.reader.Read(b)
	}
	return 0, fmt.Errorf("not opended")
}

func (f *resFile) Seek(offset int64, whence int) (int64, error) {
	if f.reader != nil {
		return f.reader.Seek(offset, whence)
	}
	return 0, fmt.Errorf("not opended")
}

func (f *resFile) Readdir(cnt int) ([]os.FileInfo, error) {
	n := len(f.item.Children)
	if cnt < n {
		n = cnt
	}
	var ret []os.FileInfo
	for _, v := range f.item.Children {
		if n <= len(ret) {
			break
		}
		ret = append(ret, &resFileStat{
			item: v,
		})
	}
	return ret, nil
}

func (f *resFile) Stat() (os.FileInfo, error) {
	return &resFileStat{
		item: f.item,
	}, nil
}

type resFileSystem struct {
	root *resfile.Item
}

func (fs *resFileSystem) Open(p string) (http.File, error) {
	paths := strings.Split(p, "/")
	if strings.HasPrefix(p, "/") {
		// remove top elements
		paths = paths[1:]
	}
	item := fs.root
	for _, s := range paths {
		if child, ok := item.Children[s]; !ok {
			return nil, fmt.Errorf("file:%q not found", s)
		} else {
			item = child
		}
	}
	f := &resFile{
		item: item,
	}
	if !item.IsDir {
		b, err := getResourceFile(item.ID)
		if err != nil {
			return nil, err
		}
		f.reader = bytes.NewReader(b)
	}
	return f, nil
}

func (d *windows) ResourcesFileSystem() (http.FileSystem, error) {
	if resRootItem == nil {
		if err := initResourceFileSystem(); err != nil {
			return nil, err
		}
	}
	return &resFileSystem{
		root: resRootItem,
	}, nil
}
