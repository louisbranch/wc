package http

import (
	"errors"
	"net/http"

	"github.com/larissavoigt/wildcare"
	"github.com/larissavoigt/wildcare/internal/view"
)

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if r.Method != "GET" || path != "/" {
		view.NotFound(w)
		return
	}

	content := struct {
		User *wildcare.User
	}{}

	u, ok := CurrentUser(r.Context())

	if ok {
		content.User = u
	}

	view.Render(w, "home/index", content)
}

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	_, ok := CurrentUser(r.Context())

	if ok {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	content := struct{ Error error }{}

	switch r.Method {
	case "GET":
		view.Render(w, "user/signup", content)
	case "POST":
		r.ParseForm()
		email := r.Form.Get("email")
		password := r.Form.Get("password")

		user := &wildcare.User{
			Email:    email,
			Password: password,
		}

		err := h.UserService.Create(user)
		if err != nil {
			content.Error = err
			view.Render(w, "user/signup", content)
			return
		}

		_, err = h.CreateSession(w, user.ID)
		if err != nil {
			content.Error = err
			view.Render(w, "user/signup", content)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	default:
		view.NotFound(w)
	}
}
func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	_, ok := CurrentUser(r.Context())

	if ok {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	content := struct{ Error error }{}

	switch r.Method {
	case "GET":
		view.Render(w, "user/login", content)
	case "POST":
		r.ParseForm()
		email := r.Form.Get("email")
		password := r.Form.Get("password")

		u, ok := h.AuthenticateUser(email, password)

		if !ok {
			content.Error = errors.New("Email and password combination doesn't match")
			view.Render(w, "user/login", content)
			return
		}

		_, err := h.CreateSession(w, u.ID)
		if err != nil {
			content.Error = errors.New("Failed to create session")
			view.Render(w, "user/login", content)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	default:
		view.NotFound(w)
	}
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	h.DestroySession(w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}
