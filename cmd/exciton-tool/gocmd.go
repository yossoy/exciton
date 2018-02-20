package main

import (
	"os/exec"
	"strings"
)

func goBuild(te *targetEnv, src string, env []string, args ...string) error {
	return goCmd(te, "build", []string{src}, env, args...)
}
func goInstall(te *targetEnv, srcs []string, env []string, args ...string) error {
	return goCmd(te, "install", srcs, env, args...)
}

func goCmd(te *targetEnv, subcmd string, srcs []string, env []string, args ...string) error {
	cmd := exec.Command(
		"go",
		subcmd,
		//"-pkgdir="+pkgdir(env), //TODO: local package for non host env
	)
	bt := te.ctx.BuildTags
	if flagBuildRelease {
		bt = append(bt, "release")
	}
	if len(bt) > 0 {
		cmd.Args = append(cmd.Args, "-tags", strings.Join(bt, " "))
	}
	if te.he.verbose {
		cmd.Args = append(cmd.Args, "-v")
	}
	if subcmd != "install" && flagBuildI {
		cmd.Args = append(cmd.Args, "-i")
	}
	if flagBuildX {
		cmd.Args = append(cmd.Args, "-x")
	}
	if flagBuildForceRebuild {
		cmd.Args = append(cmd.Args, "-a")
	}
	if flagBuildGcflags != "" {
		cmd.Args = append(cmd.Args, "-gcflags", flagBuildGcflags)
	}

	var ldflags []string
	for _, f := range strings.Fields(flagBuildLdflags) {
		if strings.HasPrefix(f, "-") {
			ldflags = append(ldflags, f)
		} else {
			ldflags[len(ldflags)-1] = ldflags[len(ldflags)-1] + " " + f
		}
	}
	var addLdflags []string
	if len(te.env.GoLdFlags) > 0 {
		addLdflags = append(addLdflags, te.env.GoLdFlags...)
	}
	if flagBuildRelease {
		addLdflags = append(addLdflags, "-w", "-s")
	}
	if len(addLdflags) > 0 {
		for _, f := range ldflags {
			for idx, ff := range addLdflags {
				if f == ff {
					addLdflags[idx] = ""
				}
			}
		}
		for _, f := range addLdflags {
			if f != "" {
				ldflags = append(ldflags, f)
			}
		}
	}

	if len(ldflags) > 0 {
		cmd.Args = append(cmd.Args, "-ldflags", strings.Join(ldflags, " "))
	}
	if flagBuildWork {
		cmd.Args = append(cmd.Args, "-work")
	}
	cmd.Args = append(cmd.Args, args...)
	cmd.Args = append(cmd.Args, srcs...)
	cmd.Env = append([]string{}, env...)
	return runCmd(te.he, cmd)
}
