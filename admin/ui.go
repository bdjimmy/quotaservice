// Licensed under the Apache License, Version 2.0
// Details: https://raw.githubusercontent.com/maniksurtani/quotaservice/master/LICENSE

package admin

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/maniksurtani/quotaservice/config"
	"github.com/maniksurtani/quotaservice/logging"
)

// uiHandler is an http.Handler for the web UI.
type uiHandler struct {
	a               Administrable
	templates       *template.Template
	htmlFiles       []string
	assetsDirectory string
	development     bool
}

func NewUIHandler(admin Administrable, assetsDirectory string, development bool) (u *uiHandler) {
	h := &uiHandler{
		a:               admin,
		assetsDirectory: assetsDirectory,
		development:     development}

	if err := h.loadTemplates(); err != nil {
		panic(err)
	}

	return h
}

func (h *uiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.development {
		if err := h.loadTemplates(); err != nil {
			logging.Printf("Caught error %v reloading templates", err)
		}
	}

	tpl := strings.TrimPrefix(r.URL.Path, "/admin/")

	if tpl == "" {
		tpl = "index.html"
	}

	err := h.templates.ExecuteTemplate(w, tpl, h.a.Configs())
	if err != nil {
		logging.Printf("Caught error %v serving URL %v", err, r.URL.Path)
		http.NotFound(w, r)
	}
}

func (h *uiHandler) loadTemplates() error {
	files, err := ioutil.ReadDir(h.assetsDirectory)

	if err != nil {
		return err
	}

	h.htmlFiles = make([]string, 0)
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".html") {
			h.htmlFiles = append(h.htmlFiles, h.assetsDirectory+"/"+f.Name())
		}
	}

	h.templates = template.Must(template.New("admin").Funcs(template.FuncMap{
		"FQN": config.FQN}).ParseFiles(h.htmlFiles...))

	return nil
}
