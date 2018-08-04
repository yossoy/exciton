// +build nonative,web

package main

import "github.com/yossoy/exciton/driver/web"

func main() {
	if err := web.Startup(ExcitonStartup); err != nil {
		panic(err)
	}
}
