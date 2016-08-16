package http

import (
	"errors"
	"net/http"

	"github.com/larissavoigt/wildcare"
	"github.com/larissavoigt/wildcare/internal/http/view"
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

	u, ok := currentUser(r.Context())

	if ok {
		content.User = u
	}

	view.Render(w, "home/index", content)
}

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	_, ok := currentUser(r.Context())

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

		user := &wildcare.User{Email: email}

		err := h.AuthenticationService.HashPassword(user, password)
		if err != nil {
			content.Error = err
			view.Render(w, "user/signup", content)
			return
		}

		err = h.UserService.Create(user)
		if err != nil {
			content.Error = err
			view.Render(w, "user/signup", content)
			return
		}

		_, err = h.createSession(w, user.ID)
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
	_, ok := currentUser(r.Context())

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

		u, err := h.UserService.FindByEmail(email)

		if err != nil {
			content.Error = errors.New("Email and password combination doesn't match")
			view.Render(w, "user/login", content)
			return
		}

		ok := h.AuthenticationService.AuthenticateUser(u, password)

		if !ok {
			content.Error = errors.New("Email and password combination doesn't match")
			view.Render(w, "user/login", content)
			return
		}

		_, err = h.createSession(w, u.ID)
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
	h.destroySession(w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}
