package window

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/yossoy/exciton/assets"
	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/log"
	"github.com/yossoy/exciton/markup"
)

type htmlContext struct {
	ID                    string
	Title                 string
	Lang                  string
	DriverType            string
	ResourcesURL          template.URL
	MesonJS               template.JS
	NativeRequestJSMethod template.JS
	JS                    []string
	CSS                   []string
	ComponentCSSFiles     []template.URL
	ComponentJSFiles      []template.URL
	IsReleaseBuild        bool
}

func loadFromAssets(fn string) ([]byte, error) {
	f, err := assets.FileSystem.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func toTemplateURL(ss []string) []template.URL {
	var r []template.URL
	for _, s := range ss {
		r = append(r, template.URL(s))
	}
	return r
}

func rootHTMLHandler(w http.ResponseWriter, r *http.Request) {
	log.PrintDebug("rootHTMLHandler: %q", r.RequestURI)
	vars := driver.RequestVars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, http.ErrNoLocation.Error(), http.StatusNotFound)
		return
	}
	win := getWindowByID(id)
	if win == nil {
		http.Error(w, http.ErrNoLocation.Error(), http.StatusNotFound)
		return
	}
	if win.cachedHTML == nil {
		a, err := loadFromAssets("/default.gohtml")
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		rfs, err := driver.ResourcesFileSystem()
		var css []string
		if err == nil {
			if f, err := rfs.Open("/css/"); err == nil {
				if fis, err := f.Readdir(-1); err == nil {
					for _, fi := range fis {
						css = append(css, fi.Name())
					}
				}
				f.Close()
			}
		}
		ctx := htmlContext{
			ID:         id,
			Title:      win.title,
			Lang:       win.lang,
			DriverType: driver.Type(),
			CSS:        css,
			NativeRequestJSMethod: template.JS(driver.NativeRequestJSMethod()),
			IsReleaseBuild:        driver.ReleaseBuild,
			ComponentCSSFiles:     toTemplateURL(markup.GetComponentCSSURLs()),
			ComponentJSFiles:      toTemplateURL(markup.GetComponentJSURLs()),
		}
		t, err := template.New("").Parse(string(a))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		var b bytes.Buffer
		if err = t.Execute(&b, ctx); err != nil {
			http.Error(w, err.Error(), http.StatusProcessing)
			return
		}
		log.PrintDebug("%s\n", b.String())
		win.cachedHTML = b.Bytes()
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(win.cachedHTML)
}

func initHTML(info *driver.StartupInfo) error {
	markup.HandleComponentResource(info.Router)
	info.Router.HandleFunc(info.AppURLBase+"/window/{id}/", rootHTMLHandler)
	//TODO: assetsのマウントは別の場所で行う?
	info.Router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(assets.FileSystem)))
	return nil
}
