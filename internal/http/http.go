package http

import (
	"context"
	"net/http"

	"github.com/larissavoigt/wildcare"
)

type Handler struct {
	UserService    wildcare.UserService
	SessionService wildcare.SessionService
}

func (h *Handler) ListenAndServe(addr string) error {
	http.Handle("/signup", h.userMiddleware(h.signup))
	http.Handle("/login", h.userMiddleware(h.login))
	http.HandleFunc("/logout", h.logout)
	http.Handle("/", h.userMiddleware(h.index))

	return http.ListenAndServe(addr, nil)
}

func (h *Handler) newUserContext(r *http.Request) context.Context {
	ctx := r.Context()

	cookie, err := r.Cookie("session")

	if err != nil {
		return ctx
	}

	s, err := h.SessionService.Find(cookie.Value)

	if err != nil {
		return ctx
	}

	u, err := h.UserService.Find(s.UserID)
	if err != nil {
		return ctx
	}

	return context.WithValue(ctx, "user", u)
}

func currentUser(ctx context.Context) (*wildcare.User, bool) {
	u, ok := ctx.Value("user").(*wildcare.User)
	return u, ok
}

func (h *Handler) userMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := h.newUserContext(r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
