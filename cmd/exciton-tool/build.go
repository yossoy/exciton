package main

import (
	"bufio"
	"errors"
	"fmt"
	"go/build"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/yossoy/exciton/markup"
)

type targetEnv struct {
	ctx build.Context
	pkg *build.Package
	env *BuildEnv
	he  *hostEnv
	//target buildTarget
	//arch   string
}

func (te *targetEnv) BuildWorkDir() string {
	return filepath.Join(te.he.tmpdir, "Work-"+te.env.OS+"-"+te.env.Arch)
}

var cmdBuild = &command{
	run:   runBuild,
	Name:  "build",
	Usage: "[build flags] [package]",
	Short: "compile exciton app",
	Long: `
	Build compiles and encodes the app named by the import path.
`,
}

// "Build flags", used by multiple commands.
var (
	flagBuildForceRebuild bool     // -a
	flagBuildI            bool     // -i
	flagBuildN            bool     // -n
	flagBuildV            bool     // -v
	flagBuildX            bool     // -x
	flagBuildO            string   // -o
	flagBuildGcflags      string   // -gcflags
	flagBuildLdflags      string   // -ldflags
	flagBuildTarget       []string // -target
	flagBuildWork         bool     // -work
	flagBuildTags         []string // -tags
	flagBuildWorkDir      string   // -w
	flagBuildRelease      bool     // -release
)

//TODO: app has vendering exciton case
var nmRE = regexp.MustCompile(`[0-9a-f]{8} t (github.com/yossoy/exciton/[^.]*)`)

func extractPkgs(he *hostEnv, nm string, path string) (map[string]bool, error) {
	if he.noExec {
		return map[string]bool{"github.com/yossoy/exciton": true}, nil
	}
	r, w := io.Pipe()
	cmd := exec.Command(nm, path)
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	nmpkgs := make(map[string]bool)
	errc := make(chan error, 1)
	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			if res := nmRE.FindStringSubmatch(s.Text()); res != nil {
				//exclude vendor directory?
				nmpkgs[res[1]] = true
			}
		}
		errc <- s.Err()
	}()

	err := cmd.Run()
	w.Close()
	if err != nil {
		return nil, fmt.Errorf("%s %s: %v", nm, path, err)
	}
	if err := <-errc; err != nil {
		return nil, fmt.Errorf("%s %s: %v", nm, path, err)
	}
	return nmpkgs, nil
}

func addBuildFlags(cmd *command) {
	cmd.flag.StringVar(&flagBuildO, "o", "", "output filename")
	cmd.flag.StringVar(&flagBuildGcflags, "gcflags", "", "arguments to pass on each go tool compile invocation.")
	cmd.flag.StringVar(&flagBuildLdflags, "ldflags", "", "arguments to pass on each go tool link invocation.")
	cmd.flag.StringVar(&flagBuildWorkDir, "w", "", "specify working directory.")

	cmd.flag.BoolVar(&flagBuildForceRebuild, "a", false, "force rebuilding of packages that are already up-to-date.")
	cmd.flag.BoolVar(&flagBuildI, "i", false, "???")                       //TODO: no need this option?
	cmd.flag.BoolVar(&flagBuildRelease, "release", false, "release build") //TODO: option => command?

	cmd.flag.Var((*stringsFlag)(&flagBuildTarget), "target", "a space-separated list of build target.")
	cmd.flag.Var((*stringsFlag)(&flagBuildTags), "tags", `a space-separated list of build tags to consider satisfied during the build.
	For more information about build tags, see the description of build constraints
	in the documentation for the go/build package.`)
}

func addBuildFlagsNVXWork(cmd *command) {
	cmd.flag.BoolVar(&flagBuildN, "n", false, "print the commands but do not run them.")
	cmd.flag.BoolVar(&flagBuildV, "v", false, "print the names of packages as they are compiled.")
	cmd.flag.BoolVar(&flagBuildX, "x", false, "print the commands.")
	cmd.flag.BoolVar(&flagBuildWork, "work", false, "print the name of the temporary work directory and do not delete it when exiting.")
}

func parseBuildTarget() (ret []*buildTargetArch, err error) {
	// for current os target
	if len(flagBuildTarget) == 0 {
		for bt := buildTarget(0); bt < buildTargetMax; bt++ {
			if bt.OSName() == runtime.GOOS {
				for _, arch := range bt.archList() {
					if arch == runtime.GOARCH {
						ret = append(ret, &buildTargetArch{target: bt, arch: arch})
					}
				}
			}
		}
		return
	}

	// for all target
	if findInSlice(flagBuildTarget, "all") >= 0 {
		//TODO: need to change ouput file (or folder)
		for bt := buildTarget(0); bt < buildTargetMax; bt++ {
			archs := bt.archList()
			for _, arch := range archs {
				ret = append(ret, &buildTargetArch{target: bt, arch: arch})
			}
		}
		return
	}
	targets := make(map[string][]string)

	for _, t := range flagBuildTarget {
		ta := strings.Split(t, "-")
		if len(ta) > 2 {
			return nil, fmt.Errorf("invalid target: %q", t)
		}
		if aa, ok := targets[ta[0]]; ok {
			if len(ta) == 1 {
				targets[ta[0]] = nil
			} else {
				if aa != nil {
					targets[ta[0]] = append(aa, ta[1])
				}
			}
		} else {
			if len(ta) == 1 {
				targets[ta[0]] = nil
			} else {
				targets[ta[0]] = []string{ta[1]}
			}
		}
	}
	for bt := buildTarget(0); bt < buildTargetMax; bt++ {
		tarchs, ok := targets[bt.String()]
		if !ok {
			continue
		}
		archs := bt.archList()
		for _, arch := range archs {
			found := false
			if tarchs == nil {
				found = true
			} else {
				for _, a := range tarchs {
					if arch == a {
						found = true
						break
					}
				}
			}
			if found {
				ret = append(ret, &buildTargetArch{target: bt, arch: arch})
			}
		}
	}
	if len(ret) == 0 {
		return nil, fmt.Errorf("invalid target: %q", flagBuildTarget)
	}
	return
}

func collectPackageResourceFiles(te *targetEnv, resDstPath string) error {
	for _, s := range te.pkg.Imports {
		pctx := build.Default
		p2, err := pctx.Import(s, te.he.cwd, build.ImportComment)
		if err == nil {
			importMarkup := false
			for _, pp := range p2.Imports {
				if strings.HasSuffix(pp, "github.com/yossoy/exciton/markup") {
					importMarkup = true
					break
				}
			}
			if !importMarkup {
				continue
			}
			resFolder := filepath.Join(p2.Dir, "resources")
			fi, err := os.Stat(resFolder)
			if os.IsNotExist(err) {
				continue
			}
			if !fi.IsDir() {
				continue
			}

			pkgFilePath := filepath.FromSlash(p2.ImportPath)
			err = filepath.Walk(resFolder, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if name := filepath.Base(path); strings.HasPrefix(name, ".") {
					// Do not include the hidden files.
					return nil
				}
				if info.IsDir() {
					return nil
				}
				relSrcPath := path[len(resFolder)+1:]

				dstFolder := filepath.Join(resDstPath, pkgFilePath, filepath.Dir(relSrcPath))
				if _, err := os.Stat(dstFolder); os.IsNotExist(err) {
					mkdir(te.he, dstFolder)
				}
				dstFilePath := filepath.Join(dstFolder, filepath.Base(relSrcPath))
				if filepath.Ext(path) == ".gocss" || filepath.Ext(path) == ".gojs" {
					b, err := markup.ReadComponentNamespaceFile(".", path, pkgFilePath)
					if err != nil {
						return err
					}
					return writeFile(te.he, dstFilePath, func(w io.Writer) error {
						_, err := w.Write(b)
						return err
					})
				}
				return copyFile(te.he, dstFilePath, path)
			})
		}
	}
	return nil
}

func runBuildOne(he *hostEnv, bta *buildTargetArch, cmd *command) error {
	args := cmd.flag.Args()
	ctxt := build.Default
	ctxt.GOARCH = bta.arch
	ctxt.GOOS = bta.target.OSName()
	var pkg *build.Package
	var err error

	switch len(args) {
	case 0:
		pkg, err = ctxt.ImportDir(he.cwd, build.ImportComment)
	case 1:
		pkg, err = ctxt.Import(args[0], he.cwd, build.ImportComment)
	default:
		cmd.usage()
		os.Exit(1)
	}
	if err != nil {
		return err
	}
	be, err := makeBuildEnv(bta.target, bta.arch)
	if err != nil {
		return err
	}
	te := &targetEnv{
		ctx: ctxt,
		pkg: pkg,
		env: be,
		he:  he,
	}
	if pkg.Name != "main" {
		return errors.New("required main package")
	}

	var nmpkgs map[string]bool
	switch bta.target {
	case buildTargetOSX:
		// TODO: use targetArchs?
		if !xcodeAvailable() {
			return fmt.Errorf("-target=osx requires XCode")
		}
		nmpkgs, err = goDarwinBuild(te, flagBuildO)
		if err != nil {
			return err
		}
	case buildTargetWindows:
		if err := gccAvailable(te.env); err != nil {
			return fmt.Errorf("-target=windows requires gcc/g++/nm: %q", err)
		}
		nmpkgs, err = goWindowsBuild(te, flagBuildO)
		if err != nil {
			return err
		}
	}
	fmt.Printf("%#v\n", nmpkgs)

	return errors.New("not implement yet")
}

func runBuild(he *hostEnv, cmd *command) error {
	//args := cmd.flag.Args()

	targets, err := parseBuildTarget()
	if err != nil {
		return fmt.Errorf(`invalid -target=%q: %v`, flagBuildTarget, err)
	}

	for _, target := range targets {
		err = runBuildOne(he, target, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	addBuildFlags(cmdBuild)
	addBuildFlagsNVXWork(cmdBuild)

}
