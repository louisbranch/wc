package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/larissavoigt/wildcare/internal/db"
	"github.com/larissavoigt/wildcare/internal/user"
)

type Session struct {
	Token   string
	ID      int64
	Expires time.Time
}

func Create(w http.ResponseWriter, id int64) (*Session, error) {
	t, err := token()

	if err != nil {
		return nil, err
	}

	s := &Session{
		Token:   t,
		ID:      id,
		Expires: time.Now().Add(7 * 24 * time.Hour),
	}

	_, err = db.Exec("INSERT INTO sessions (token, user_id, expires_at) VALUES(?, ?, ?)",
		s.Token, s.ID, s.Expires)

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

func NewContext(r *http.Request) context.Context {
	ctx := r.Context()

	token, err := r.Cookie("session")

	if err != nil {
		return ctx
	}

	var id int64

	err = db.QueryRow(`
		SELECT user_id
		FROM sessions where token=?`, token.Value).Scan(&id)

	if err != nil {
		return ctx
	}

	u, err := user.Find(id)
	if err != nil {
		return ctx
	}

	return context.WithValue(ctx, "user", u)
}

func CurrentUser(ctx context.Context) (*user.User, bool) {
	u, ok := ctx.Value("user").(*user.User)
	return u, ok
}

func Destroy(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("session")
	if err != nil {
		return
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	_, err = db.Exec("DELETE FROM sessions where token = ? LIMIT 1", token.Value)

	if err != nil {
		log.Printf("[WARN] session record not found: %s", err)
	}
}

func Middleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
