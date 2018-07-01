package main

import (
	"fmt"
	"runtime"
)

var cmdTargets = &command{
	run:   runListTargets,
	Name:  "targets",
	Usage: "",
	Short: "list available build targets",
	Long: `
	List available build targets.
`,
}

func runListTargets(he *hostEnv, cmd *command) error {
	for bt := buildTarget(0); bt < buildTargetMax; bt++ {
		archs := bt.archList()
		for _, arch := range archs {
			if bt.OSName() == runtime.GOOS && arch == runtime.GOARCH {
				fmt.Printf("* ")
			} else {
				fmt.Printf("- ")
			}
			fmt.Printf("%s-%s\n", bt.String(), arch)
		}
	}
	return nil
}
