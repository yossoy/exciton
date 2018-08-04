// +build !nonative

package main

import (
	"github.com/yossoy/exciton/driver/windows"
)

func main() {
	if err := windows.Startup(ExcitonStartup); err != nil {
		panic(err)
	}
}
