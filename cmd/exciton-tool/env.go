package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type hostEnv struct {
	tmpdir       string
	cleanupFuncs []func()
	xout         io.Writer
	verbose      bool
	noExec       bool
	preserveWork bool
	hostOS       string
	cwd          string
}

func (be *hostEnv) finalize() {
	for _, cf := range be.cleanupFuncs {
		cf()
	}
}

func initHostEnv() (*hostEnv, error) {
	be := &hostEnv{}
	if flagBuildN || flagBuildV {
		be.verbose = true
	}
	if flagBuildWork {
		be.preserveWork = true
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	be.cwd = cwd
	if flagBuildN {
		be.noExec = true
		be.tmpdir = "$WORK"
	} else if flagBuildWorkDir != "" {
		be.noExec = false
		be.tmpdir = flagBuildWorkDir
	} else {
		tmpdir, err := ioutil.TempDir("", "exciton-work-")
		if err != nil {
			return nil, err
		}
		be.tmpdir = tmpdir
		be.cleanupFuncs = append(be.cleanupFuncs, func() {
			removeAll(be, tmpdir)
		})
	}
	be.xout = os.Stderr

	return be, nil
}

// environ merges os.Environ and the given "key=value" pairs.
// If a key is in both os.Environ and kv, kv takes precedence.
func environ(be *hostEnv, kv []string) []string {
	cur := os.Environ()
	new := make([]string, 0, len(cur)+len(kv))

	envs := make(map[string]string, len(cur))
	for _, ev := range cur {
		elem := strings.SplitN(ev, "=", 2)
		if len(elem) != 2 || elem[0] == "" {
			// pass the env var of unusual form untouched.
			// e.g. Windows may have env var names starting with "=".
			new = append(new, ev)
			continue
		}
		if be.hostOS == "windows" {
			elem[0] = strings.ToUpper(elem[0])
		}
		envs[elem[0]] = elem[1]
	}
	for _, ev := range kv {
		elem := strings.SplitN(ev, "=", 2)
		if len(elem) != 2 || elem[0] == "" {
			panic(fmt.Sprintf("malformed env var %q from input", ev))
		}
		if be.hostOS == "windows" {
			elem[0] = strings.ToUpper(elem[0])
		}
		envs[elem[0]] = elem[1]
	}
	for k, v := range envs {
		new = append(new, k+"="+v)
	}
	return new
}
