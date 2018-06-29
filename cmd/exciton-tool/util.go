package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func goEnv(name string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}
	val, err := exec.Command("go", "env", name).Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(val))
}

func printcmd(be *hostEnv, format string, args ...interface{}) {
	cmd := fmt.Sprintf(format+"\n", args...)
	if be.tmpdir != "" {
		cmd = strings.Replace(cmd, be.tmpdir, "$WORK", -1)
	}
	if androidHome := os.Getenv("ANDROID_HOME"); androidHome != "" {
		cmd = strings.Replace(cmd, androidHome, "$ANDROID_HOME", -1)
	}
	// if gomobilepath != "" {
	// 	cmd = strings.Replace(cmd, gomobilepath, "$GOMOBILE", -1)
	// }
	if goroot := goEnv("GOROOT"); goroot != "" {
		cmd = strings.Replace(cmd, goroot, "$GOROOT", -1)
	}
	if gopath := goEnv("GOPATH"); gopath != "" {
		cmd = strings.Replace(cmd, gopath, "$GOPATH", -1)
	}
	if env := os.Getenv("HOME"); env != "" {
		cmd = strings.Replace(cmd, env, "$HOME", -1)
	}
	if env := os.Getenv("HOMEPATH"); env != "" {
		cmd = strings.Replace(cmd, env, "$HOMEPATH", -1)
	}
	fmt.Fprint(be.xout, cmd)
}

func runCmd(be *hostEnv, cmd *exec.Cmd) error {
	if be.verbose {
		dir := ""
		if cmd.Dir != "" {
			dir = "PWD=" + cmd.Dir + " "
		}
		env := strings.Join(cmd.Env, " ")
		if env != "" {
			env += " "
		}
		printcmd(be, "%s%s%s", dir, env, strings.Join(cmd.Args, " "))
	}

	buf := new(bytes.Buffer)
	buf.WriteByte('\n')
	if be.verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = buf
		cmd.Stderr = buf
	}

	if be.preserveWork {
		if runtime.GOOS == "windows" {
			cmd.Env = append(cmd.Env, `TEMP=`+be.tmpdir)
			cmd.Env = append(cmd.Env, `TMP=`+be.tmpdir)
		} else {
			cmd.Env = append(cmd.Env, `TMPDIR=`+be.tmpdir)
		}
	}

	if !be.noExec {
		cmd.Env = environ(be, cmd.Env)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s failed: %v%s", strings.Join(cmd.Args, " "), err, buf)
		}
	}
	return nil
}

func removeAll(be *hostEnv, path string) error {
	if be.verbose {
		printcmd(be, `rm -r -f "%s"`, path)
	}
	if be.noExec {
		return nil
	}

	// os.RemoveAll behaves differently in windows.
	// http://golang.org/issues/9606
	if be.hostOS == "windows" {
		resetReadOnlyFlagAll(path)
	}

	return os.RemoveAll(path)
}

func resetReadOnlyFlagAll(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return os.Chmod(path, 0666)
	}
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	names, _ := fd.Readdirnames(-1)
	for _, name := range names {
		resetReadOnlyFlagAll(path + string(filepath.Separator) + name)
	}
	return nil
}

func mkdir(he *hostEnv, dir string) error {
	if he.verbose {
		printcmd(he, "mkdir -p %s", dir)
	}
	if he.noExec {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}

func copyFile(he *hostEnv, dst, src string) error {
	if he.verbose {
		printcmd(he, "cp %s %s", src, dst)
	}
	return writeFile(he, dst, func(w io.Writer) error {
		if he.noExec {
			return nil
		}
		f, err := os.Open(src)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(w, f); err != nil {
			return fmt.Errorf("cp %s %s failed: %v", src, dst, err)
		}
		return nil
	})
}

func writeFile(he *hostEnv, filename string, generate func(io.Writer) error) error {
	if he.verbose {
		fmt.Fprintf(os.Stderr, "write %s\n", filename)
	}

	err := mkdir(he, filepath.Dir(filename))
	if err != nil {
		return err
	}

	if he.noExec {
		return generate(ioutil.Discard)
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); err == nil {
			err = cerr
		}
	}()

	return generate(f)
}

func findInSlice(array []string, str string) int {
	for i, s := range array {
		if s == str {
			return i
		}
	}
	return -1
}
