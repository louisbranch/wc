package view

import (
	"fmt"
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

func New(names ...string) *View {
	for i, n := range names {
		p := []string{"templates"}
		p = append(p, strings.Split(n+".html", "/")...)
		names[i] = filepath.Join(p...)
	}
	return &View{paths: names}
}

func (v View) Render(w http.ResponseWriter, data interface{}) {
	for k, v := range views.cache {
		fmt.Printf("%s: %s\n", k, v.Name())
	}

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
	v := New("layouts/base", path)
	v.Render(w, data)
}

func NotFound(w http.ResponseWriter) {
	Render(w, "errors/404", nil)
}

func parse(names ...string) (tpl *template.Template, err error) {
	cp := make([]string, len(names))
	copy(cp, names)
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
