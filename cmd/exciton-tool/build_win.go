package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

const manifestFile = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
 <dependency>
   <dependentAssembly>
     <assemblyIdentity
       type="win32"
       name="Microsoft.Windows.Common-Controls"
       version="6.0.0.0"
       publicKeyToken="6595b64144ccf1df"
       language="*"
       processorArchitecture="*"/>
   </dependentAssembly>
 </dependency>
 <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
   <security>
     <requestedPrivileges>
       <requestedExecutionLevel
         level="asInvoker"
         uiAccess="false"/>
       </requestedPrivileges>
   </security>
 </trustInfo>
</assembly>
`

//TODO: versions as template
const rcTempl = `#include "winuser.h"
1			RT_MANIFEST	{{.ManifestName}}

{{if .IconPath}}IDI_ICON1	ICON		{{.IconPath}}{{end}}

1			VERSIONINFO
FILEVERSION		1,0,0,0
PRODUCTVERSION	1,0,0,0
FILEFLAGS		0x0L
FILEFLAGSMASK	0x3fL
FILESUBTYPE		0
BEGIN
	BLOCK	"StringFileInfo"
	BEGIN
		BLOCK	"000004b0"
		BEGIN
			VALUE "CompanyName",      ""
			VALUE "FileDescription",  "{{.FileDescription}}"
			VALUE "FileVersion",      "1.0"
			VALUE "InternalName",     "{{.Name}}"
			VALUE "LegalCopyright",   "{{.Copyright}}"
			VALUE "OriginalFilename", "{{.Name}}"
			VALUE "ProductName",      "{{.ProductName}}"
			VALUE "ProductVersion",   "1.0.0.0"
		END
	END
	BLOCK "VarFileInfo"
    BEGIN
        VALUE "Translation", 0x0, 1200 // Neutral language, Unicode
	END
END
`

func generateWinResourceFile(te *targetEnv, outFile string, resFilePath string) error {
	pkg := te.pkg
	src := pkg.ImportPath
	workdir := te.BuildWorkDir()
	srcManifestName := filepath.Join(src, filepath.Base(src)+".exe.manifest")
	manFile, err := ioutil.TempFile(workdir, "manifest")
	if err != nil {
		return err
	}
	defer os.Remove(manFile.Name())
	defer manFile.Close()
	if _, err := os.Stat(srcManifestName); err == nil {
		// exist original name
		err = func() error {
			src, err := os.Open(srcManifestName)
			if err != nil {
				return err
			}
			defer src.Close()
			_, err = io.Copy(manFile, src)

			return err
		}()
		if err != nil {
			return err
		}
	} else {
		_, err = manFile.WriteString(manifestFile)
		if err != nil {
			return err
		}
	}
	manFile.Close()

	rcFile, err := ioutil.TempFile(workdir, "rc")
	if err != nil {
		return err
	}
	defer os.Remove(rcFile.Name())
	defer rcFile.Close()
	tmpl, err := template.New("resTemplate").Parse(rcTempl)
	if err != nil {
		return err
	}

	args := struct {
		ManifestName    string
		IconPath        string
		FileDescription string
		Name            string
		Copyright       string
		ProductName     string
	}{
		ManifestName: filepath.Base(manFile.Name()),
	}
	args.Name = outFile + ".exe"
	args.ProductName = outFile

	err = tmpl.Execute(rcFile, &args)
	if err != nil {
		return err
	}

	rcFile.Close()

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(pwd)
	err = os.Chdir(workdir)
	if err != nil {
		return err
	}

	cmdarg := []string{
		fmt.Sprintf("--input=%s", filepath.Base(rcFile.Name())),
		fmt.Sprintf("--output=%s", resFilePath),
		"--output-format=coff",
	}

	cmd := exec.Command(te.env.WINDRES, cmdarg...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
func goWindowsBuild(te *targetEnv, outFile string) (map[string]bool, error) {
	pkg := te.pkg
	ctx := te.ctx
	be := te.env
	workdir := te.BuildWorkDir()
	src := pkg.ImportPath
	if strings.HasSuffix(outFile, ".exe") {
		return nil, fmt.Errorf("-o must suppress an .exe")
	}

	mkdir(te.he, workdir)

	envStrs := be.envStrings()
	ctx.BuildTags = append(ctx.BuildTags, be.Target.String())

	if outFile == "" {
		outFile = path.Base(pkg.ImportPath)
	}

	var args []string
	var targetPath string
	if flagBuildRelease {
		mkdir(te.he, filepath.Join(workdir, "output"))
		targetPath = filepath.Join(workdir, "output", outFile+".exe")
		args = append(args, "-o="+targetPath)
	} else {
		targetPath = outFile + ".exe"
	}

	// generate syso if need
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	resFileName := filepath.Join(cwd, "resources.syso")
	if _, err := os.Stat(resFileName); err != nil {
		err := generateWinResourceFile(te, outFile, resFileName)
		if err != nil {
			return nil, err
		}
		defer os.Remove(resFileName)
	}

	//TODO: リソースを作る

	//TODO: アウトプットのフォルダはどうする?
	//ターゲットと自分が一致する場合は同じ場所で良いだろうけど。。。
	//あと、テンポラリフォルダも素のまま使っても良いという訳ではないかと
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

	// copy resources
	copyFile(te.he, filepath.Join(te.he.cwd, path.Base(targetPath)), targetPath)

	return nmpkgs, nil
}
