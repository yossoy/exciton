package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type BuildEnv struct {
	Target           buildTarget
	OS               string
	Arch             string
	CompilerSpecific map[string]string
	CC               string
	CXX              string
	NM               string
	WINDRES          string
	CFLAGS           string
	CXXFLAGS         string
	LDFLAGS          string
	CgoEnable        bool
	GoLdFlags        []string
}

func (be *BuildEnv) envStrings() []string {
	var ret []string
	ret = append(ret, "GOOS="+be.OS)
	ret = append(ret, "GOARCH="+be.Arch)
	for k, v := range be.CompilerSpecific {
		ret = append(ret, k+"="+v)
	}
	ret = append(ret, "CC="+be.CC)
	ret = append(ret, "CXX="+be.CXX)
	ret = append(ret, "CGO_CFLAGS="+be.CFLAGS)
	ret = append(ret, "CGO_CXXFLAGS="+be.CXXFLAGS)
	ret = append(ret, "CGO_LDFLAGS="+be.LDFLAGS)
	if be.CgoEnable {
		ret = append(ret, "CGO_ENABLED=1")
	}
	return ret
}

func archClang(arch string) string {
	switch arch {
	case "arm":
		return "armv7"
	case "arm64":
		return "arm64"
	case "386":
		return "i386"
	case "amd64":
		return "x86_64"
	default:
		panic("unknown arch: " + arch)
	}
}

func envClang(sdkName string) (clang, cflags string, err error) {
	cmd := exec.Command("xcrun", "--sdk", sdkName, "--find", "clang")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("xcrun --find: %v\n%s", err, out)
	}
	clang = strings.TrimSpace(string(out))

	cmd = exec.Command("xcrun", "--sdk", sdkName, "--show-sdk-path")
	out, err = cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("xcrun --show-sdk-path: %v\n%s", err, out)
	}
	sdk := strings.TrimSpace(string(out))

	return clang, "-isysroot " + sdk, nil
}

func makeDarwinBuildEnv(target buildTarget, arch string) (*BuildEnv, error) {
	clang, cflags, err := envClang(target.darwinSdkName(arch))
	if err != nil {
		return nil, err
	}
	be := &BuildEnv{
		Target:    target,
		OS:        target.OSName(),
		Arch:      arch,
		CC:        clang,
		CXX:       clang,
		NM:        "nm",
		CFLAGS:    cflags + " -arch " + archClang(arch),
		LDFLAGS:   cflags + " -arch " + archClang(arch), //TODO: minversion
		CgoEnable: true,
	}
	if arch == "arm" {
		be.CompilerSpecific = map[string]string{"GOARM": "7"}
	}
	verKey, minVer := target.sdkRequireVersionVersion()
	if verKey != "" && minVer != "" {
		be.CFLAGS += " -m" + verKey + "=" + minVer
	}
	be.CXXFLAGS = be.CFLAGS

	return be, nil
}

func makeGccBuildEnv(target buildTarget, arch string) (*BuildEnv, error) {
	var gccMachine string
	switch arch {
	case "386":
		gccMachine = "i686"
	case "amd64":
		gccMachine = "x86_64"
	default:
		return nil, fmt.Errorf("Unsupported arch: %q", arch)
	}
	var gccOS string
	var gccVendor string
	switch target.OSName() {
	case "windows":
		gccVendor = "w64"
		gccOS = "mingw32"
	default:
		return nil, fmt.Errorf("Unsupported OS: %q", target.OSName())
	}
	var gccTriplet string
	if gccVendor == "" {
		gccTriplet = gccMachine + "-" + gccOS
	} else {
		gccTriplet = gccMachine + "-" + gccVendor + "-" + gccOS
	}
	be := &BuildEnv{
		Target:    target,
		OS:        target.OSName(),
		Arch:      arch,
		CC:        gccTriplet + "-gcc",
		CXX:       gccTriplet + "-g++",
		NM:        gccTriplet + "-nm",
		CgoEnable: true,
	}
	if target == buildTargetWindows {
		be.WINDRES = gccTriplet + "-windres"
		if flagBuildRelease {
			be.GoLdFlags = []string{"-H windowsgui"}
		}
	}

	//Test compiler exists?

	return be, nil
}

func makeBuildEnv(target buildTarget, arch string) (*BuildEnv, error) {
	switch target {
	case buildTargetOSX:
		return makeDarwinBuildEnv(target, arch)
	case buildTargetWindows:
		return makeGccBuildEnv(target, arch)
	default:
		return nil, errors.New("unsupported buildTarget: " + target.String())
	}
}
