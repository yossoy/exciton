package markup

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/yossoy/exciton/log"

	"github.com/gorilla/mux"
)

func HandleComponentResource(r *mux.Router) {
	r.HandleFunc("/components/{id}/resources/{filename}", componentResourceFileHandle)
}

func escapeKlassPathForURL(s string) string {
	var b strings.Builder
	for _, c := range s {
		switch c {
		case '/':
			b.WriteString(fmt.Sprintf("_%02x", c))
		case '%':
			b.WriteString("__")
		default:
			b.WriteRune(c)
		}

	}
	return b.String()
}

func unescapeKlassPathForURL(s string) (string, error) {
	var b strings.Builder
	state := 0
	var prevVal int
	for _, c := range s {
		switch state {
		case 0:
			switch c {
			case '_':
				state = 1
			default:
				b.WriteRune(c)
			}
		case 1:
			switch {
			case '0' <= c && c <= '9':
				prevVal = int(c - '0')
				state = 2
			case 'a' <= c && c <= 'f':
				prevVal = int(10 + (c - 'a'))
				state = 2
			case 'A' <= c && c <= 'F':
				prevVal = int(10 + (c - 'A'))
				state = 2
			case c == '_':
				b.WriteRune(c)
			default:
				return "", fmt.Errorf("invalid source: %s", s)
			}
		case 2:
			switch {
			case '0' <= c && c <= '9':
				b.WriteRune(rune(prevVal*16 + int(c-'0')))
				state = 0
			case 'a' <= c && c <= 'f':
				b.WriteRune(rune(prevVal*16 + 10 + int(c-'a')))
				state = 0
			case 'A' <= c && c <= 'F':
				b.WriteRune(rune(prevVal*16 + 10 + int(c-'A')))
				state = 0
			default:
				return "", fmt.Errorf("invalid source: %s", s)
			}
		}
	}
	return b.String(), nil
}

func componentResourceFileHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, http.ErrNoLocation.Error(), http.StatusNotFound)
		return
	}
	klassName, err := unescapeKlassPathForURL(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	log.PrintDebug("Unescaped: %q\n", klassName)
	k, ok := componentKlasses[klassName]
	if !ok {
		http.Error(w, http.ErrNoLocation.Error(), http.StatusNotFound)
		return
	}
	log.PrintDebug("Found Klass: %q\n", klassName)
	fn, ok := vars["filename"]
	if !ok {
		http.Error(w, http.ErrNoLocation.Error(), http.StatusNotFound)
		return
	}
	f, err := k.GetResourceFile(fn)
	if err != nil {
		log.PrintDebug("Not found file: %q", fn)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	// ext := path.Ext(fn)
	w.Header().Add("Date", fi.ModTime().Format(http.TimeFormat))
	w.Header().Add("Content-Type", mime.TypeByExtension(path.Ext(fn)))
	// switch ext {
	// case ".css":
	// }

	w.Header().Add("Content-Length", fmt.Sprintf("%d", fi.Size()))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, f)
}

func GetComponentCSSURLs() []string {
	var cssURLs []string
	for _, k := range componentKlasses {
		if k.cssFile != "" {
			p := fmt.Sprintf("/components/%s/resources/%s", escapeKlassPathForURL(k.Name()), k.cssFile)
			cssURLs = append(cssURLs, p)
		}
	}
	return cssURLs
}

func GetComponentJSURLs() []string {
	var jsURLs []string
	for _, k := range componentKlasses {
		if k.jsFile != "" {
			p := fmt.Sprintf("/components/%s/resources/%s", escapeKlassPathForURL(k.Name()), k.jsFile)
			jsURLs = append(jsURLs, p)
		}
	}
	return jsURLs
}
