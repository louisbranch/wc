package user

import (
	"github.com/larissavoigt/wildcare/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID    int64
	Name  string
	Email string
}

func Create(email, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return nil, err
	}

	user := &User{Email: email}

	result, err := db.Exec(`
	INSERT INTO users (email, password_hash)
	VALUES(?, ?)
	`, user.Email, hash)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return nil, err
	}

	user.ID = id

	return user, nil
}

func Find(id int64) (*User, error) {
	u := &User{}
	err := db.QueryRow(`
		SELECT id, name, email
		FROM users where id=?`, id).Scan(
		&u.ID, &u.Name, &u.Email)

	return u, err
}

func Authenticate(email, password string) (*User, bool) {
	u := &User{}
	var hash string
	err := db.QueryRow(`
		SELECT id, name, email, password_hash
		FROM users where email=?`, email).Scan(
		&u.ID, &u.Name, &u.Email, &hash)
	if err != nil {
		return nil, false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return u, err == nil
}
