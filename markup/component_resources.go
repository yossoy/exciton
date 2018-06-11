package markup

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/mux"
)

func HandleComponentResource(r *mux.Router) {
	r.HandleFunc("/components/{id}/{name}/resources/{filename}", componentResourceFileHandle)
}

func getKlassFromIDandName(id, name string) *Klass {
	if kpi, ok := klassPathIDs[id]; ok {
		if k, ok := kpi.klasses[name]; ok {
			return k
		}
	}
	return nil
}

func componentResourceFileHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, http.ErrNoLocation.Error(), http.StatusNotFound)
		return
	}
	klassName, ok := vars["name"]
	if !ok {
		http.Error(w, http.ErrNoLocation.Error(), http.StatusNotFound)
		return
	}

	k := getKlassFromIDandName(id, klassName)
	if k == nil {
		http.Error(w, http.ErrNoLocation.Error(), http.StatusNotFound)
		return
	}
	fn, ok := vars["filename"]
	if !ok {
		http.Error(w, http.ErrNoLocation.Error(), http.StatusNotFound)
		return
	}
	f, err := k.GetResourceFile(fn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	ext := path.Ext(fn)
	w.Header().Add("Date", fi.ModTime().Format(http.TimeFormat))
	w.Header().Add("Content-Type", mime.TypeByExtension(path.Ext(fn)))
	if !strings.HasSuffix(path.Base(fn), "-global") {
		switch ext {
		case ".css":
			css, err := ioutil.ReadAll(f)
			if err != nil {
				http.Error(w, err.Error(), http.StatusProcessing)
				return
			}
			clsPrefix := id + "-" + strings.TrimSuffix(fn, ext)
			s, _, err := convertKlassCSS(string(css), clsPrefix)
			if err != nil {
				http.Error(w, err.Error(), http.StatusProcessing)
				return
			}
			w.Header().Add("Content-Length", fmt.Sprintf("%d", len(s)))
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, s)
			return
		default:
			break
		}
	}

	w.Header().Add("Content-Length", fmt.Sprintf("%d", fi.Size()))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, f)
}

func GetComponentCSSURLs() []string {
	var cssURLs []string
	for kpid, kpi := range klassPathIDs {
		for _, k := range kpi.klasses {
			if k.cssFile != "" {
				p := fmt.Sprintf("/components/%s/%s/resources/%s", kpid, k.name, k.cssFile)
				cssURLs = append(cssURLs, p)
			}
		}
	}
	return cssURLs
}

func GetComponentJSURLs() []string {
	var jsURLs []string
	for kpid, kpi := range klassPathIDs {
		for _, k := range kpi.klasses {
			if k.jsFile != "" {
				p := fmt.Sprintf("/components/%s/%s/resources/%s", kpid, k.name, k.jsFile)
				jsURLs = append(jsURLs, p)
			}
		}
	}
	return jsURLs
}
