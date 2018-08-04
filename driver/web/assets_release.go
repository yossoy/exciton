// +build release

package web

import (
	"net/http"
)

var fileSystem http.FileSystem = assetsData
