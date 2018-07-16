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

	"github.com/yossoy/exciton/driver"

	"github.com/gorilla/mux"
)

const (
	jsTemplate = `(function() {
var id = %q;
function wrapper(exports, require, module, __filename, __dirname) {
	%s
};
console.log('register', id);
window.exciton.registerModule(id, wrapper);
})();`
)

func (k *Klass) ResourcePathBase() string {
	return fmt.Sprintf("/components/%s", k.pathInfo.id)
}

func HandleComponentResource(r driver.Router) {
	r.PathPrefix("/components/{id}/resources/").HandlerFunc(componentResourceFileHandle)
}

func getKlassPathInfoFromID(id string) *klassPathInfo {
	if kpi, ok := klassPathIDs[id]; ok {
		return kpi
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

	kpi := getKlassPathInfoFromID(id)
	if kpi == nil {
		http.Error(w, http.ErrNoLocation.Error(), http.StatusNotFound)
		return
	}
	f, err := kpi.getResourceFile(fn)
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
			if g, ok := kpi.cssFiles[fn]; !ok || g {
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
		case ".js":
			if g, ok := kpi.jsFiles[fn]; !ok || g {
				break
			}
			js, err := ioutil.ReadAll(f)
			if err != nil {
				http.Error(w, err.Error(), http.StatusProcessing)
				return
			}
			js = stripShebang(js)
			prefix := id + "-" + strings.TrimSuffix(fn, ext)
			jsStr := fmt.Sprintf(jsTemplate, prefix, string(js))
			fmt.Printf("*********** jsStr = %s\n", jsStr)
			w.Header().Add("Content-Length", fmt.Sprintf("%d", len(jsStr)))
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, jsStr)
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
		for k := range kpi.cssFiles {
			p := fmt.Sprintf("/components/%s/resources/%s", kpid, k)
			cssURLs = append(cssURLs, p)
		}
	}
	return cssURLs
}

func GetComponentJSURLs() []string {
	var jsURLs []string
	for kpid, kpi := range klassPathIDs {
		for k := range kpi.jsFiles {
			p := fmt.Sprintf("/components/%s/resources/%s", kpid, k)
			jsURLs = append(jsURLs, p)
		}
	}
	return jsURLs
}

func stripShebang(content []byte) []byte {
	if len(content) < 2 {
		return content
	}
	if content[0] != '#' || content[1] != '!' {
		return content
	}
	if len(content) == 2 {
		return content[:0]
	}
	var i int
	for i = 2; i < len(content); i++ {
		if content[i] == '\r' || content[i] == '\n' {
			break
		}
	}
	if i == len(content) {
		return content[:0]
	}
	return content[i:]
}
