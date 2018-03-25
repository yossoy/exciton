package window

import (
	"bytes"
	"html/template"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/yossoy/exciton/markup"

	"github.com/yossoy/exciton/assets"
	"github.com/yossoy/exciton/driver"
)

type htmlContext struct {
	ID                    string
	Title                 string
	Lang                  string
	IsIE                  bool
	ResourcesURL          template.URL
	MesonJS               template.JS
	NativeRequestJSMethod template.JS
	JS                    []string
	CSS                   []string
	ComponentCSS          []template.CSS
	ComponentCSSFiles     []template.URL
	ComponentJS           []template.JS
	ComponentJSFiles      []template.URL
	IsReleaseBuild        bool
}

func setupHTML(cfg *WindowConfig) error {
	js, err := assets.Asset("assets/exciton.js")
	if err != nil {
		return err
	}
	// make resources folder's file url
	resPath, err := driver.Resources()
	if err != nil {
		return err
	}
	resources := filepath.ToSlash(resPath)
	if !strings.HasPrefix(resources, "/") {
		resources = "/" + resources
	}
	if !strings.HasSuffix(resources, "/") {
		resources = resources + "/"
	}
	resurl := &url.URL{}
	resurl.Scheme = "file"
	resurl.Host = ""
	resurl.Path = resources

	ctx := htmlContext{
		ID:                    cfg.ID,
		Title:                 cfg.Title,
		Lang:                  cfg.Lang,
		IsIE:                  driver.IsIE(),
		ResourcesURL:          template.URL(resurl.String()),
		MesonJS:               template.JS(string(js)),
		NativeRequestJSMethod: template.JS(driver.NativeRequestJSMethod()),
		IsReleaseBuild:        releaseBuild,
	}

	csss, err := markup.GetComponentCSSFiles(resPath)
	if err != nil {
		return err
	}
	jss, err := markup.GetComponentJSFiles(resPath)
	if err != nil {
		return err
	}
	if releaseBuild {
		ctx.ComponentCSSFiles = make([]template.URL, len(csss))
		for i, p := range csss {
			ctx.ComponentCSSFiles[i] = template.URL(p)
		}
		ctx.ComponentJSFiles = make([]template.URL, len(jss))
		for i, p := range jss {
			ctx.ComponentJSFiles[i] = template.URL(p)
		}
	} else {
		ctx.ComponentCSS = make([]template.CSS, len(csss))
		for i, css := range csss {
			ctx.ComponentCSS[i] = template.CSS(css)
		}
		ctx.ComponentJS = make([]template.JS, len(jss))
		for i, js := range jss {
			ctx.ComponentJS[i] = template.JS(js)
		}
	}

	a, err := assets.Asset("assets/default.gohtml")
	if err != nil {
		return err
	}
	t, err := template.New("").Parse(string(a))
	if err != nil {
		return err
	}
	var b bytes.Buffer
	if err = t.Execute(&b, ctx); err != nil {
		return err
	}
	cfg.HTML = b.String()
	return nil
}
