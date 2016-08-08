package main

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/larissavoigt/wildcare/internal/session"
	"github.com/larissavoigt/wildcare/internal/user"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").ParseGlob("templates/*.html"))
}

func main() {

	http.Handle("/", session.Middleware(home))
	http.Handle("/signup", session.Middleware(signup))
	http.Handle("/login", session.Middleware(login))

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		session.Destroy(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		tpl.ExecuteTemplate(w, "404", nil)
		return
	}

	content := struct {
		User *user.User
	}{}

	u, ok := session.CurrentUser(r.Context())

	if ok {
		content.User = u
	}

	tpl.ExecuteTemplate(w, "index.html", content)
}

func signup(w http.ResponseWriter, r *http.Request) {
	_, ok := session.CurrentUser(r.Context())

	if ok {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	content := struct{ Error error }{}

	switch r.Method {
	case "GET":
		tpl.ExecuteTemplate(w, "signup.html", content)
	case "POST":
		r.ParseForm()
		email := r.Form.Get("email")
		password := r.Form.Get("password")

		user, err := user.Create(email, password)
		if err != nil {
			content.Error = err
			tpl.ExecuteTemplate(w, "signup.html", content)
			return
		}

		_, err = session.Create(w, user.ID)
		if err != nil {
			content.Error = err
			tpl.ExecuteTemplate(w, "signup.html", content)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	default:
		tpl.ExecuteTemplate(w, "404", nil)
	}
}
func login(w http.ResponseWriter, r *http.Request) {
	_, ok := session.CurrentUser(r.Context())

	if ok {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	content := struct{ Error error }{}

	switch r.Method {
	case "GET":
		tpl.ExecuteTemplate(w, "login.html", content)
	case "POST":
		r.ParseForm()
		email := r.Form.Get("email")
		password := r.Form.Get("password")

		u, ok := user.Authenticate(email, password)

		if !ok {
			content.Error = errors.New("Email and password combination doesn't match")
			tpl.ExecuteTemplate(w, "login.html", content)
			return
		}

		_, err := session.Create(w, u.ID)
		if err != nil {
			content.Error = errors.New("Failed to create session")
			tpl.ExecuteTemplate(w, "login.html", content)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	default:
		tpl.ExecuteTemplate(w, "404", nil)
	}
}
