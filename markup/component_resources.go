package markup

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/yossoy/exciton/log"

	"github.com/gorilla/mux"
)

func (k *Klass) ResourcePathBase() string {
	return fmt.Sprintf("/components/%s/%s", k.pathInfo.id, k.name)
}

func HandleComponentResource(r *mux.Router) {
	r.PathPrefix("/components/{id}/{name}/resources/").HandlerFunc(componentResourceFileHandle)
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
	route := mux.CurrentRoute(r)
	reg, _ := route.GetPathRegexp()
	rc := regexp.MustCompile(reg)
	fs := rc.FindString(r.URL.String())
	fn := strings.TrimPrefix(r.URL.String(), fs)
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
		log.PrintDebug("getKlassFromIDandName class not found\n")
		http.Error(w, http.ErrNoLocation.Error(), http.StatusNotFound)
		return
	}
	f, err := k.GetResourceFile(fn)
	if err != nil {
		log.PrintDebug("k.GetResourceFile() failed\n = %v\n", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		log.PrintDebug("k.GetResourceFile() failed\n = %v\n", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	ext := path.Ext(fn)
	w.Header().Add("Date", fi.ModTime().Format(http.TimeFormat))
	w.Header().Add("Content-Type", mime.TypeByExtension(path.Ext(fn)))
	if !strings.HasSuffix(path.Base(fn), "-global") {
		switch ext {
		case ".css":
			if k.cssIsGlobal {
				break
			}
			css, err := ioutil.ReadAll(f)
			if err != nil {
				http.Error(w, err.Error(), http.StatusProcessing)
				return
			}
			clsPrefix := id + "-" + escapeClassName(strings.TrimSuffix(fn, ext))
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
