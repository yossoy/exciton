package main

import (
	"github.com/yossoy/exciton/driver/mac"
)

func main() {
	if err := mac.Startup(ExcitonStartup); err != nil {
		panic(err)
	}
}
