package wildcare

import "time"

type User struct {
	ID       int64
	Name     string
	Email    string
	Password string
}

type Session struct {
	Token   string
	UserID  int64
	Expires time.Time
}

type UserService interface {
	Authenticate(email, password string) (*User, bool)
	Find(id int64) (*User, error)
	Create(*User) error
}

type SessionService interface {
	Create(*Session) error
	Delete(*Session) error
	Find(token string) (*Session, error)
}
