// +build release

package assets

import (
	"net/http"
)

var FileSytem http.FileSystem = assetsData
