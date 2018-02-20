package main

import (
	"os/exec"
)

func gccAvailable(be *BuildEnv) error {
	//TODO: check gcc versions?
	if _, err := exec.LookPath(be.CC); err != nil {
		return err
	}
	if _, err := exec.LookPath(be.CXX); err != nil {
		return err
	}
	if _, err := exec.LookPath(be.NM); err != nil {
		return err
	}
	return nil
}
