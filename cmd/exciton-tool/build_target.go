package main

import "fmt"

type buildTargetArch struct {
	target buildTarget
	arch   string
}

type buildTarget int

const (
	buildTargetOSX buildTarget = iota
	buildTargetWindows
	// buildTargetLinux
	// buildTargetIPhone
	// buildTargetIPhoneSim
	// buildTargetAndroid
	buildTargetMax
)

//const goosList = "android darwin dragonfly freebsd linux nacl netbsd openbsd plan9 solaris windows zos "
//const goarchList = "386 amd64 amd64p32 arm armbe arm64 arm64be ppc64 ppc64le mips mipsle mips64 mips64le mips64p32 mips64p32le ppc s390 s390x sparc sparc64 "

func (bt buildTarget) OSName() string {
	switch bt {
	case buildTargetOSX:
		return "darwin"
	case buildTargetWindows:
		return "windows"
	default:
		panic("unsupported buildTarget")
	}
}

func (bt buildTarget) String() string {
	switch bt {
	case buildTargetOSX:
		return "osx"
	// case buildTargetIPhone:
	// 	return "ios"
	case buildTargetWindows:
		return "windows"
	default:
		return fmt.Sprintf("Unknown Target: %d", bt)

	}

}

func (bt buildTarget) archList() []string {
	switch bt {
	case buildTargetOSX:
		return []string{"amd64"} // dropped 386
	case buildTargetWindows:
		return []string{"386", "amd64"}
	default:
		panic("unsupported buildTarget" + bt.String())
	}
}

//TODO: move build_darwin?
func (bt buildTarget) darwinSdkName(arch string) string {
	switch bt {
	case buildTargetOSX:
		return "macosx"
	// case buildTargetIPhone:
	// 	switch arch {
	// 	case "386", "amd64":
	// 		return "iphonesimulator"
	// 	case "arm", "arm64":
	// 		return "iphoneos"
	// 	default:
	// 		panic("invalid darwin arch: " + arch)
	// 	}
	default:
		panic("invalid darwin target: " + bt.String())
	}
}

func (bt buildTarget) sdkRequireVersionVersion() (string, string) {
	switch bt {
	case buildTargetOSX:
		return "macosx-version-min", "10.10" //TODO:
	// case buildTargetIPhone:
	// 	return "9.0"
	default:
		return "", ""
	}
}
