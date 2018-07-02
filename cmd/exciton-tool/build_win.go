package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/yossoy/exciton/driver/windows/resfile"
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

{{range .AppendResFiles}}{{.ID}}	RCDATA 	"{{.Path}}"
{{end}}
`

func addWin32ResFileItems(id int, items resfile.Items, paths []string, srcDir string) (string, error) {
	if len(paths) == 0 {
		return "", nil
	}
	p := paths[0]
	paths = paths[1:]
	srcDir = filepath.Join(srcDir, p)
	fi, err := os.Stat(srcDir)
	if err != nil {
		return "", err
	}
	item, ok := items[p]
	if !ok {
		item = &resfile.Item{
			ID:       0,
			Name:     p,
			Size:     fi.Size(),
			Mode:     fi.Mode(),
			ModTime:  fi.ModTime(),
			IsDir:    fi.IsDir(),
			Children: make(resfile.Items),
		}
		items[p] = item
	}
	if !fi.IsDir() {
		if len(paths) == 0 {
			item.ID = id
			return srcDir, nil
		}
		return "", fmt.Errorf("invalid path %q(file) + %q", srcDir, filepath.Join(paths...))
	}

	return addWin32ResFileItems(id, item.Children, paths, srcDir)
}

type idAndFile struct {
	ID   int
	Path string
}

func toWin32ResFileItem(files []*collectFileItem) (resfile.Items, []idAndFile, error) {
	items := make(resfile.Items)
	var resFiles []idAndFile
	//fmt.Println("*** Res Files ***")
	id := resfile.FileIDStart
	for _, file := range files {
		rootItem := items
		if file.dstRelDir != "" {
			for _, ss := range strings.Split(file.dstRelDir, string(filepath.Separator)) {
				item, ok := rootItem[ss]
				if !ok {
					item = &resfile.Item{
						ID:       0,
						Name:     ss,
						Size:     0,
						Mode:     0,
						ModTime:  time.Now(),
						IsDir:    true,
						Children: make(resfile.Items),
					}
					rootItem[ss] = item
				}
				rootItem = item.Children
			}
		}
		//fmt.Printf("toWind32ResFileItem: %q\n", file.dstRelDir)
		for _, f := range file.files {
			paths := strings.Split(f, string(os.PathSeparator))
			fp, err := addWin32ResFileItems(id, rootItem, paths, file.srcDir)
			if err != nil {
				return nil, nil, err
			}
			if fp != "" {
				resFiles = append(resFiles, idAndFile{
					ID:   id,
					Path: filepath.ToSlash(fp),
				})
				id++
			}
		}
	}
	if len(items) == 0 {
		if len(resFiles) != 0 {
			panic(false)
		}
		return nil, nil, nil
	}
	return items, resFiles, nil
}

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

	resFiles, err := collectPackageResourceFileItems(te)
	if err != nil {
		return err
	}
	resItems, fileList, err := toWin32ResFileItem(resFiles)
	if err != nil {
		return err
	}
	//jsonResItem, err := json.Marshal(resItems)
	jsonResItem, err := json.MarshalIndent(resItems, "", "  ")
	if err != nil {
		return err
	}
	jsonResFile, err := ioutil.TempFile(workdir, "jsonResItem")
	if n, err := jsonResFile.Write(jsonResItem); err != nil {
		return err
	} else if n != len(jsonResItem) {
		panic(false)
	}
	jsonResFile.Close()
	absJSONResPath, err := filepath.Abs(jsonResFile.Name())
	if err != nil {
		return err
	}
	fileList = append(fileList, idAndFile{
		ID:   resfile.FileMapJsonID,
		Path: filepath.ToSlash(absJSONResPath),
	})

	// fmt.Printf("json:\n%s\n", string(jsonResItem))
	// for _, itm := range fileList {
	// 	fmt.Printf("[%d] %q\n", itm.ID, itm.Path)
	// }

	args := struct {
		ManifestName    string
		IconPath        string
		FileDescription string
		Name            string
		Copyright       string
		ProductName     string
		AppendResFiles  []idAndFile
	}{
		ManifestName:   filepath.Base(manFile.Name()),
		AppendResFiles: fileList,
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
