package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	ttmpl "text/template"
)

func goDarwinBuild(te *targetEnv, outFile string) (map[string]bool, error) {
	pkg := te.pkg
	ctx := te.ctx
	be := te.env
	src := pkg.ImportPath
	workdir := te.BuildWorkDir()
	if flagBuildO != "" && strings.HasSuffix(outFile, ".app") {
		return nil, fmt.Errorf("-o must suppress an .app")
	}

	envStrs := be.envStrings()
	ctx.BuildTags = append(ctx.BuildTags, be.Target.String())

	var args []string
	var targetPath string
	if flagBuildRelease {
		targetPath = filepath.Join(workdir, be.Arch)
		args = append(args, "-o="+targetPath)
	} else {
		targetPath = path.Base(pkg.ImportPath)
	}
	if err := goBuild(te, src, envStrs, args...); err != nil {
		return nil, err
	}
	nmpkgs, err := extractPkgs(te.he, be.NM, targetPath)
	if err != nil {
		return nil, err
	}
	if !flagBuildRelease {
		return nmpkgs, nil
	}

	productName := rfc1034Label(path.Base(pkg.ImportPath))
	if productName == "" {
		productName = "ProductName" // like xcode.
	}

	infoplist := new(bytes.Buffer)
	if err := infoplistTmpl.Execute(infoplist, infoplistTmplData{
		// TODO: better bundle id. present config or reversed package path?
		BundleID: "org.golang.todo." + productName,
		Name:     strings.Title(path.Base(pkg.ImportPath)),
		OSName:   be.Target.OSName(),
	}); err != nil {
		return nil, err
	}

	projPbxproj := new(bytes.Buffer)
	if err := projPbxprojTmpl.Execute(projPbxproj, projPbxprojData{
		SdkName: be.Target.darwinSdkName(be.Arch),
	}); err != nil {
		return nil, err
	}

	files := []struct {
		name     string
		contents []byte
	}{
		{filepath.Join(workdir, "main.xcodeproj", "project.pbxproj"), projPbxproj.Bytes()},
		{filepath.Join(workdir, "main", "Info.plist"), infoplist.Bytes()},
		{filepath.Join(workdir, "main", "Images.xcassets", "AppIcon.appiconset", "Contents.json"), []byte(contentsJSON)},
	}

	for _, file := range files {
		if err := mkdir(te.he, filepath.Dir(file.name)); err != nil {
			return nil, err
		}
		if te.he.verbose {
			printcmd(te.he, "echo \"%s\" > %s", file.contents, file.name)
		}
		if !te.he.noExec {
			if err := ioutil.WriteFile(file.name, file.contents, 0644); err != nil {
				return nil, err
			}
		}
	}

	// arm64Path := filepath.Join(tmpdir, "arm64")
	// if err := goBuild(src, darwinArm64Env, "-o="+arm64Path); err != nil {
	// 	return nil, err
	// }

	// Apple requires builds to target both darwin/arm and darwin/arm64.
	// We are using lipo tool to build multiarchitecture binaries.
	// TODO(jbd): Investigate the new announcements about iO9's fat binary
	// size limitations are breaking this feature.
	//TODO: no need lipo?
	cmd := exec.Command(
		"xcrun", "lipo",
		"-create", targetPath,
		"-o", filepath.Join(workdir, "main", "main"),
	)
	if err := runCmd(te.he, cmd); err != nil {
		return nil, err
	}

	// TODO(jbd): Set the launcher icon.
	if err := darwinCopyAssets(te, workdir); err != nil {
		return nil, err
	}

	// Build and move the release build to the output directory.
	cmd = exec.Command(
		"xcrun", "xcodebuild",
		"-configuration", "Release",
		"-project", filepath.Join(workdir, "main.xcodeproj"),
	)
	if err := runCmd(te.he, cmd); err != nil {
		return nil, err
	}

	// TODO(jbd): Fallback to copying if renaming fails.
	if outFile == "" {
		n := pkg.ImportPath
		if n == "." {
			// use cwd name
			cwd, err := os.Getwd()
			if err != nil {
				return nil, fmt.Errorf("cannot create .app; cannot get the current working dir: %v", err)
			}
			n = cwd
		}
		n = path.Base(n)
		outFile = n
	}
	var outFolderSuffix string
	switch be.Target.OSName() {
	case "ios":
		outFolderSuffix = "-iphoneos"
	default:
		outFolderSuffix = ""
	}
	outFile = outFile + ".app"
	if te.he.verbose {
		printcmd(te.he, "mv %s %s", filepath.Join(workdir, "build", "Release"+outFolderSuffix, "main.app"), outFile)
	}
	if !te.he.noExec {
		// if output already exists, remove.
		if err := os.RemoveAll(outFile); err != nil {
			return nil, err
		}
		if err := os.Rename(filepath.Join(workdir, "build", "Release"+outFolderSuffix, "main.app"), outFile); err != nil {
			return nil, err
		}
	}
	return nmpkgs, nil
}

func xcodeAvailable() bool {
	_, err := exec.LookPath("xcrun")
	return err == nil
}

func darwinCopyAssets(te *targetEnv, xcodeProjDir string) error {
	he := te.he
	dstAssets := xcodeProjDir + "/main/assets"
	if err := mkdir(he, dstAssets); err != nil {
		return err
	}

	srcAssets := filepath.Join(te.pkg.Dir, "resources")
	fi, err := os.Stat(srcAssets)
	if err != nil {
		if os.IsNotExist(err) {
			// skip walking through the directory to deep copy.
			return nil
		}
		return err
	}
	if !fi.IsDir() {
		// skip walking through to deep copy.
		return nil
	}
	// if assets is a symlink, follow the symlink.
	srcAssets, err = filepath.EvalSymlinks(srcAssets)
	if err != nil {
		return err
	}
	err = filepath.Walk(srcAssets, func(path string, info os.FileInfo, err error) error {
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
		dst := dstAssets + "/" + path[len(srcAssets)+1:]
		return copyFile(he, dst, path)
	})

	collectPackageResourceFiles(te, dstAssets)

	return err
}

type infoplistTmplData struct {
	BundleID string
	Name     string
	OSName   string
}

var infoplistTmpl = ttmpl.Must(ttmpl.New("infoplist").Parse(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>BuildMachineOSBuild</key>
  <string>14F1509</string>
  <key>CFBundleDevelopmentRegion</key>
  <string>en</string>
  <key>CFBundleExecutable</key>
  <string>main</string>
  <key>CFBundleIdentifier</key>
  <string>{{.BundleID}}</string>
  <key>CFBundleInfoDictionaryVersion</key>
  <string>6.0</string>
  <key>CFBundleName</key>
  <string>{{.Name}}</string>
  <key>CFBundlePackageType</key>
  <string>APPL</string>
  <key>CFBundleShortVersionString</key>
  <string>1.0</string>
  <key>CFBundleSignature</key>
  <string>????</string>
  <key>CFBundleVersion</key>
  <string>1</string>
  {{if eq .OSName "ios"}}
  <key>LSRequiresIPhoneOS</key>
  <true/>{{end}}
  <key>UILaunchStoryboardName</key>
  <string>LaunchScreen</string>
  <key>UISupportedInterfaceOrientations</key>
  <array>
    <string>UIInterfaceOrientationPortrait</string>
    <string>UIInterfaceOrientationLandscapeLeft</string>
    <string>UIInterfaceOrientationLandscapeRight</string>
  </array>
  <key>UISupportedInterfaceOrientations~ipad</key>
  <array>
    <string>UIInterfaceOrientationPortrait</string>
    <string>UIInterfaceOrientationPortraitUpsideDown</string>
    <string>UIInterfaceOrientationLandscapeLeft</string>
    <string>UIInterfaceOrientationLandscapeRight</string>
  </array>
</dict>
</plist>
`))

type projPbxprojData struct {
	SdkName string
}

var projPbxprojTmpl = template.Must(template.New("pbxproj").Parse(`// !$*UTF8*$!
{
  archiveVersion = 1;
  classes = {
  };
  objectVersion = 46;
  objects = {

/* Begin PBXBuildFile section */
    254BB84F1B1FD08900C56DE9 /* Images.xcassets in Resources */ = {isa = PBXBuildFile; fileRef = 254BB84E1B1FD08900C56DE9 /* Images.xcassets */; };
    254BB8681B1FD16500C56DE9 /* main in Resources */ = {isa = PBXBuildFile; fileRef = 254BB8671B1FD16500C56DE9 /* main */; settings = {ATTRIBUTES = (CodeSignOnCopy, ); }; };
    25FB30331B30FDEE0005924C /* assets in Resources */ = {isa = PBXBuildFile; fileRef = 25FB30321B30FDEE0005924C /* assets */; };
/* End PBXBuildFile section */

/* Begin PBXCopyFilesBuildPhase section */
61B5FB4D09C4E9FA00B25A18 /* Copy main Executable */ = {
	isa = PBXCopyFilesBuildPhase;
	buildActionMask = 2147483647;
	dstPath = "";
	dstSubfolderSpec = 6;
	files = (
		254BB8681B1FD16500C56DE9 /* main in Resources */,
	);
	name = "Copy main Executable";
	runOnlyForDeploymentPostprocessing = 0;
};
/* End PBXCopyFilesBuildPhase section */

/* Begin PBXFileReference section */
    254BB83E1B1FD08900C56DE9 /* main.app */ = {isa = PBXFileReference; explicitFileType = wrapper.application; includeInIndex = 0; path = main.app; sourceTree = BUILT_PRODUCTS_DIR; };
    254BB8421B1FD08900C56DE9 /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
    254BB84E1B1FD08900C56DE9 /* Images.xcassets */ = {isa = PBXFileReference; lastKnownFileType = folder.assetcatalog; path = Images.xcassets; sourceTree = "<group>"; };
    254BB8671B1FD16500C56DE9 /* main */ = {isa = PBXFileReference; lastKnownFileType = "compiled.mach-o.executable"; path = main; sourceTree = "<group>"; };
    25FB30321B30FDEE0005924C /* assets */ = {isa = PBXFileReference; lastKnownFileType = folder; name = assets; path = main/assets; sourceTree = "<group>"; };
/* End PBXFileReference section */

/* Begin PBXGroup section */
    254BB8351B1FD08900C56DE9 = {
      isa = PBXGroup;
      children = (
        25FB30321B30FDEE0005924C /* assets */,
        254BB8401B1FD08900C56DE9 /* main */,
        254BB83F1B1FD08900C56DE9 /* Products */,
      );
      sourceTree = "<group>";
      usesTabs = 0;
    };
    254BB83F1B1FD08900C56DE9 /* Products */ = {
      isa = PBXGroup;
      children = (
        254BB83E1B1FD08900C56DE9 /* main.app */,
      );
      name = Products;
      sourceTree = "<group>";
    };
    254BB8401B1FD08900C56DE9 /* main */ = {
      isa = PBXGroup;
      children = (
        254BB8671B1FD16500C56DE9 /* main */,
        254BB84E1B1FD08900C56DE9 /* Images.xcassets */,
        254BB8411B1FD08900C56DE9 /* Supporting Files */,
      );
      path = main;
      sourceTree = "<group>";
    };
    254BB8411B1FD08900C56DE9 /* Supporting Files */ = {
      isa = PBXGroup;
      children = (
        254BB8421B1FD08900C56DE9 /* Info.plist */,
      );
      name = "Supporting Files";
      sourceTree = "<group>";
    };
/* End PBXGroup section */

/* Begin PBXNativeTarget section */
    254BB83D1B1FD08900C56DE9 /* main */ = {
      isa = PBXNativeTarget;
      buildConfigurationList = 254BB8611B1FD08900C56DE9 /* Build configuration list for PBXNativeTarget "main" */;
      buildPhases = (
        254BB83C1B1FD08900C56DE9 /* Resources */,
        61B5FB4D09C4E9FA00B25A18 /* Copy */,
      );
      buildRules = (
      );
      dependencies = (
      );
      name = main;
      productName = main;
      productReference = 254BB83E1B1FD08900C56DE9 /* main.app */;
      productType = "com.apple.product-type.application";
    };
/* End PBXNativeTarget section */

/* Begin PBXProject section */
    254BB8361B1FD08900C56DE9 /* Project object */ = {
      isa = PBXProject;
      attributes = {
        LastUpgradeCheck = 0630;
        ORGANIZATIONNAME = Developer;
        TargetAttributes = {
          254BB83D1B1FD08900C56DE9 = {
            CreatedOnToolsVersion = 6.3.1;
          };
        };
      };
      buildConfigurationList = 254BB8391B1FD08900C56DE9 /* Build configuration list for PBXProject "main" */;
      compatibilityVersion = "Xcode 3.2";
      developmentRegion = English;
      hasScannedForEncodings = 0;
      knownRegions = (
        en,
        Base,
      );
      mainGroup = 254BB8351B1FD08900C56DE9;
      productRefGroup = 254BB83F1B1FD08900C56DE9 /* Products */;
      projectDirPath = "";
      projectRoot = "";
      targets = (
        254BB83D1B1FD08900C56DE9 /* main */,
      );
    };
/* End PBXProject section */

/* Begin PBXResourcesBuildPhase section */
    254BB83C1B1FD08900C56DE9 /* Resources */ = {
      isa = PBXResourcesBuildPhase;
      buildActionMask = 2147483647;
      files = (
        25FB30331B30FDEE0005924C /* assets in Resources */,
        254BB84F1B1FD08900C56DE9 /* Images.xcassets in Resources */,
      );
      runOnlyForDeploymentPostprocessing = 0;
    };
/* End PBXResourcesBuildPhase section */

/* Begin XCBuildConfiguration section */
    254BB8601B1FD08900C56DE9 /* Release */ = {
      isa = XCBuildConfiguration;
      buildSettings = {
        ALWAYS_SEARCH_USER_PATHS = NO;
        CLANG_CXX_LANGUAGE_STANDARD = "gnu++0x";
        CLANG_CXX_LIBRARY = "libc++";
        CLANG_ENABLE_MODULES = YES;
        CLANG_ENABLE_OBJC_ARC = YES;
		CLANG_WARN_BLOCK_CAPTURE_AUTORELEASING = YES;
        CLANG_WARN_BOOL_CONVERSION = YES;
		CLANG_WARN_COMMA = YES;
        CLANG_WARN_CONSTANT_CONVERSION = YES;
        CLANG_WARN_DIRECT_OBJC_ISA_USAGE = YES_ERROR;
        CLANG_WARN_EMPTY_BODY = YES;
        CLANG_WARN_ENUM_CONVERSION = YES;
		CLANG_WARN_INFINITE_RECURSION = YES;
        CLANG_WARN_INT_CONVERSION = YES;
		CLANG_WARN_NON_LITERAL_NULL_CONVERSION = YES;
		CLANG_WARN_OBJC_LITERAL_CONVERSION = YES;
        CLANG_WARN_OBJC_ROOT_CLASS = YES_ERROR;
		CLANG_WARN_RANGE_LOOP_ANALYSIS = YES;
		CLANG_WARN_STRICT_PROTOTYPES = YES;
		CLANG_WARN_SUSPICIOUS_MOVE = YES;
		CLANG_WARN_UNREACHABLE_CODE = YES;
        CLANG_WARN__DUPLICATE_METHOD_MATCH = YES;
        "CODE_SIGN_IDENTITY[sdk=iphoneos*]" = "iPhone Developer";
		COPY_PHASE_STRIP = NO;
        DEBUG_INFORMATION_FORMAT = "dwarf-with-dsym";
        ENABLE_NS_ASSERTIONS = NO;
        ENABLE_STRICT_OBJC_MSGSEND = YES;
        GCC_C_LANGUAGE_STANDARD = gnu99;
        GCC_NO_COMMON_BLOCKS = YES;
        GCC_WARN_64_TO_32_BIT_CONVERSION = YES;
        GCC_WARN_ABOUT_RETURN_TYPE = YES_ERROR;
        GCC_WARN_UNDECLARED_SELECTOR = YES;
        GCC_WARN_UNINITIALIZED_AUTOS = YES_AGGRESSIVE;
        GCC_WARN_UNUSED_FUNCTION = YES;
        GCC_WARN_UNUSED_VARIABLE = YES;
        IPHONEOS_DEPLOYMENT_TARGET = 8.3;
        MTL_ENABLE_DEBUG_INFO = NO;
        SDKROOT = {{.SdkName}};
        TARGETED_DEVICE_FAMILY = "1,2";
        VALIDATE_PRODUCT = YES;
      };
      name = Release;
    };
    254BB8631B1FD08900C56DE9 /* Release */ = {
      isa = XCBuildConfiguration;
      buildSettings = {
        ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
        INFOPLIST_FILE = main/Info.plist;
        LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/Frameworks";
        PRODUCT_NAME = "$(TARGET_NAME)";
      };
      name = Release;
    };
/* End XCBuildConfiguration section */

/* Begin XCConfigurationList section */
    254BB8391B1FD08900C56DE9 /* Build configuration list for PBXProject "main" */ = {
      isa = XCConfigurationList;
      buildConfigurations = (
        254BB8601B1FD08900C56DE9 /* Release */,
      );
      defaultConfigurationIsVisible = 0;
      defaultConfigurationName = Release;
    };
    254BB8611B1FD08900C56DE9 /* Build configuration list for PBXNativeTarget "main" */ = {
      isa = XCConfigurationList;
      buildConfigurations = (
        254BB8631B1FD08900C56DE9 /* Release */,
      );
      defaultConfigurationIsVisible = 0;
      defaultConfigurationName = Release;
    };
/* End XCConfigurationList section */
  };
  rootObject = 254BB8361B1FD08900C56DE9 /* Project object */;
}
`))

const contentsJSON = `{
  "images" : [
    {
      "idiom" : "iphone",
      "size" : "29x29",
      "scale" : "2x"
    },
    {
      "idiom" : "iphone",
      "size" : "29x29",
      "scale" : "3x"
    },
    {
      "idiom" : "iphone",
      "size" : "40x40",
      "scale" : "2x"
    },
    {
      "idiom" : "iphone",
      "size" : "40x40",
      "scale" : "3x"
    },
    {
      "idiom" : "iphone",
      "size" : "60x60",
      "scale" : "2x"
    },
    {
      "idiom" : "iphone",
      "size" : "60x60",
      "scale" : "3x"
    },
    {
      "idiom" : "ipad",
      "size" : "29x29",
      "scale" : "1x"
    },
    {
      "idiom" : "ipad",
      "size" : "29x29",
      "scale" : "2x"
    },
    {
      "idiom" : "ipad",
      "size" : "40x40",
      "scale" : "1x"
    },
    {
      "idiom" : "ipad",
      "size" : "40x40",
      "scale" : "2x"
    },
    {
      "idiom" : "ipad",
      "size" : "76x76",
      "scale" : "1x"
    },
    {
      "idiom" : "ipad",
      "size" : "76x76",
      "scale" : "2x"
    }
  ],
  "info" : {
    "version" : 1,
    "author" : "xcode"
  }
}
`

// rfc1034Label sanitizes the name to be usable in a uniform type identifier.
// The sanitization is similar to xcode's rfc1034identifier macro that
// replaces illegal characters (not conforming the rfc1034 label rule) with '-'.
func rfc1034Label(name string) string {
	// * Uniform type identifier:
	//
	// According to
	// https://developer.apple.com/library/ios/documentation/FileManagement/Conceptual/understanding_utis/understand_utis_conc/understand_utis_conc.html
	//
	// A uniform type identifier is a Unicode string that usually contains characters
	// in the ASCII character set. However, only a subset of the ASCII characters are
	// permitted. You may use the Roman alphabet in upper and lower case (A–Z, a–z),
	// the digits 0 through 9, the dot (“.”), and the hyphen (“-”). This restriction
	// is based on DNS name restrictions, set forth in RFC 1035.
	//
	// Uniform type identifiers may also contain any of the Unicode characters greater
	// than U+007F.
	//
	// Note: the actual implementation of xcode does not allow some unicode characters
	// greater than U+007f. In this implementation, we just replace everything non
	// alphanumeric with "-" like the rfc1034identifier macro.
	//
	// * RFC1034 Label
	//
	// <label> ::= <letter> [ [ <ldh-str> ] <let-dig> ]
	// <ldh-str> ::= <let-dig-hyp> | <let-dig-hyp> <ldh-str>
	// <let-dig-hyp> ::= <let-dig> | "-"
	// <let-dig> ::= <letter> | <digit>
	const surrSelf = 0x10000
	begin := false

	var res []rune
	for i, r := range name {
		if r == '.' && !begin {
			continue
		}
		begin = true

		switch {
		case 'a' <= r && r <= 'z', 'A' <= r && r <= 'Z':
			res = append(res, r)
		case '0' <= r && r <= '9':
			if i == 0 {
				res = append(res, '-')
			} else {
				res = append(res, r)
			}
		default:
			if r < surrSelf {
				res = append(res, '-')
			} else {
				res = append(res, '-', '-')
			}
		}
	}
	return string(res)
}
