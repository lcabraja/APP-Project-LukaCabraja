package http

import (
	"errors"
	"github.com/lcabraja/APP-Project-LukaCabraja/log"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

type HttpManager struct {
	prefix string
	host   string
}

func NewHttpManager(prefix, host string) *HttpManager {
	hm := newHttpManager(prefix, host)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/favicon.svg", http.StatusFound)
	})

	uiDir := "./ui/build"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.If("[/] request for: %v\n", r.RequestURI)

		if r.URL.Path == "/" {
			content, err := ioutil.ReadFile(uiDir + "/200.html")
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.Write(content)
			return
		}

		http.FileServer(http.Dir(uiDir)).ServeHTTP(w, r)
	})
	return hm
}

func newHttpManager(prefix, host string) *HttpManager {
	return &HttpManager{
		prefix: prefix,
		host:   host,
	}
}

func (hm *HttpManager) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	var p string
	if strings.HasSuffix(pattern, "/") {
		p = path.Join(hm.prefix, pattern) + "/"
	} else {

		p = path.Join(hm.prefix, pattern)
	}
	log.Devf("added listener: %s\n", p)
	http.HandleFunc(p, handler)
}

func (hm *HttpManager) ListenAndServe() {
	log.If("listening on: [%s]\n", hm.host)

	err := http.ListenAndServe(hm.host, nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.E("server closed")
	} else if err != nil {
		log.Ff("error starting server: %s\n", err)
	}
}
