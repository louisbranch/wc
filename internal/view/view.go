package view

import (
	"html/template"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type View struct {
	paths []string
}

var views = struct {
	sync.RWMutex
	cache map[string]*template.Template
}{
	cache: make(map[string]*template.Template),
}

func New(name string) *View {
	return &View{paths: []string{templatePath(name)}}
}

func (v *View) Include(name string) *View {
	v.paths = append(v.paths, templatePath(name))
	return v
}

func (v *View) Render(w http.ResponseWriter, data interface{}) {
	tpl, err := parse(v.paths...)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	err = tpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func Render(w http.ResponseWriter, path string, data interface{}) {
	v := New("layouts/base").Include(path)
	v.Render(w, data)
}

func NotFound(w http.ResponseWriter) {
	Render(w, "errors/404", nil)
}

func templatePath(name string) string {
	p := []string{"templates"}
	p = append(p, strings.Split(name+".html", "/")...)

	return filepath.Join(p...)
}

func parse(names ...string) (tpl *template.Template, err error) {
	cp := make([]string, len(names))
	sort.Strings(cp)
	id := strings.Join(cp, ":")

	views.RLock()
	tpl, ok := views.cache[id]
	views.RUnlock()

	if !ok {
		tpl, err = template.ParseFiles(names...)
		if err != nil {
			return nil, err
		}
		views.Lock()
		views.cache[id] = tpl
		views.Unlock()
	}

	return tpl, nil
}
