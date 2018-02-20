package main

import (
	"flag"
	"fmt"
	"os"
)

type command struct {
	run   func(*hostEnv, *command) error
	flag  flag.FlagSet
	Name  string
	Usage string
	Short string
	Long  string
}

func (cmd *command) usage() {
	fmt.Fprintf(os.Stdout, "usage: %s %s %s\n%s\n", excitonToolName, cmd.Name, cmd.Usage, cmd.Long)
	cmd.flag.PrintDefaults()
}
