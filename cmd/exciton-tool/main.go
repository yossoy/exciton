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
		// cmdVersion,
	}
	usageTmpl = template.Must(template.New("usage").Parse(
		`exciton-tool is a tool for building and running gui apps written in Go.

To install:

	$ go get TODO...
	$ exciton-tool init

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

func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		return // succeeded at helping
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s help command\n\nToo many arguments given.\n", excitonToolName)
		os.Exit(2) // failed to help
	}

	arg := args[0]
	for _, cmd := range commands {
		if cmd.Name == arg {
			cmd.usage()
			return // succeeded at helping
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q.  Run '%s help'.\n", arg, excitonToolName)
	os.Exit(2)
}

func main() {
	excitonToolName = os.Args[0]
	flag.Usage = func() {
		printUsage(os.Stderr)
		os.Exit(2)
	}
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	for _, cmd := range commands {
		if cmd.Name == args[0] {
			cmd.flag.Usage = func() {
				cmd.usage()
				os.Exit(1)
			}
			cmd.flag.Parse(args[1:])
			he, err := initHostEnv()
			defer he.finalize()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %q\n", os.Args[0], err)
				os.Exit(1)
			}
			if err := cmd.run(he, cmd); err != nil {
				msg := err.Error()
				if msg != "" {
					fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
				}
				os.Exit(1)
			}
			return
		}
	}
	fmt.Fprintf(os.Stderr, "%s: unknown subcommand %q\nRun 'gomobile help' for usage.\n", os.Args[0], args[0])
	os.Exit(2)
}
