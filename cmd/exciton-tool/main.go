package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
)

var (
	excitonToolName = "exciton-tool"
	goVersionOut    = []byte(nil)
	commands        = []*command{
		// cmdInit,
		cmdBuild,
		cmdTargets,
		// cmdVersion,
	}
	usageTmpl = template.Must(template.New("usage").Parse(
		`exciton-tool is a tool for building and running gui apps written in Go.

To build:

	$ go get TODO...
	$ exciton-tool build

At least Go 1.9 is required.

Usage:

	exciton-tool command [arguments]

Commands:
{{range .}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}

Use 'exciton-tool help [command]' for more information about that command.
`))
)

func printUsage(w io.Writer) {
	bufw := bufio.NewWriter(w)
	if err := usageTmpl.Execute(bufw, commands); err != nil {
		panic(err)
	}
	bufw.Flush()
}

func help(args []string) int {
	if len(args) == 0 {
		printUsage(os.Stdout)
		return 0 // succeeded at helping
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s help command\n\nToo many arguments given.\n", excitonToolName)
		os.Exit(2) // failed to help
	}

	arg := args[0]
	for _, cmd := range commands {
		if cmd.Name == arg {
			cmd.usage()
			return 0 // succeeded at helping
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q.  Run '%s help'.\n", arg, excitonToolName)
	return 2
}

func doMain() int {
	excitonToolName = os.Args[0]
	flag.Usage = func() {
		printUsage(os.Stderr)
		return
	}
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
	}

	if args[0] == "help" {
		return help(args[1:])
	}

	for _, cmd := range commands {
		if cmd.Name == args[0] {
			cmd.flag.Usage = func() {
				cmd.usage()
				return
			}
			cmd.flag.Parse(args[1:])
			he, err := initHostEnv()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %q\n", os.Args[0], err)
				return 1
			}
			defer he.finalize()
			if err := cmd.run(he, cmd); err != nil {
				msg := err.Error()
				if msg != "" {
					fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
				}
				return 1
			}
			return 0
		}
	}
	fmt.Fprintf(os.Stderr, "%s: unknown subcommand %q\nRun 'gomobile help' for usage.\n", os.Args[0], args[0])
	return 2
}

func main() {
	os.Exit(doMain())
}
