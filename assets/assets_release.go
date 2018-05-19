// +build release

package assets

import (
	"net/http"
)

var FileSystem http.FileSystem = assetsData
