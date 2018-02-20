package window

import (
	"bytes"
	"html/template"
	"net/url"
	"path/filepath"
	"strings"

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
}

func setupHTML(cfg *WindowConfig) error {
	js, err := assets.Asset("assets/exciton.js")
	if err != nil {
		return err
	}
	// make resources folder's file url
	resources, err := driver.Resources()
	if err != nil {
		return err
	}
	resources = filepath.ToSlash(resources)
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
