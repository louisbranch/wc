package http

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/larissavoigt/wildcare"
)

func (h *Handler) CreateSession(w http.ResponseWriter, id int64) (*wildcare.Session, error) {
	t, err := token()

	if err != nil {
		return nil, err
	}

	s := &wildcare.Session{
		Token:   t,
		UserID:  id,
		Expires: time.Now().Add(7 * 24 * time.Hour),
	}

	err = h.SessionService.Create(s)

	cookie := &http.Cookie{
		Name:     "session",
		Value:    s.Token,
		Path:     "/",
		Expires:  s.Expires,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	return s, nil
}

func token() (string, error) {
	size := 40
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func (h *Handler) DestroySession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return
	}

	s := &wildcare.Session{
		Token: cookie.Value,
	}

	err = h.SessionService.Delete(s)

	if err != nil {
		log.Printf("[WARN] session record not found: %s", err)
	}

	cookie = &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}

func (h *Handler) AuthenticateUser(email, password string) (*wildcare.User, bool) {
	return nil, false
}
