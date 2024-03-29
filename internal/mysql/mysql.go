package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/larissavoigt/wildcare"
)

func Open(name string) (*sql.DB, error) {
	return sql.Open("mysql", "root:@/"+name+"?parseTime=true")
}

type UserService struct {
	*sql.DB
}

type SessionService struct {
	*sql.DB
}

func (s *UserService) Create(u *wildcare.User) error {
	result, err := s.DB.Exec(`
	INSERT INTO users (email, password_hash)
	VALUES(?, ?)
	`, u.Email, u.PasswordHash)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	u.ID = id

	return nil
}

func (db *UserService) Find(id int64) (*wildcare.User, error) {
	u := &wildcare.User{}
	err := db.QueryRow(`
		SELECT id, name, email
		FROM users where id=?`, id).Scan(
		&u.ID, &u.Name, &u.Email)

	return u, err
}

func (db *UserService) FindByEmail(email string) (*wildcare.User, error) {
	u := &wildcare.User{}
	err := db.QueryRow(`
		SELECT id, name, email, password_hash
		FROM users where email=?`, email).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash)

	return u, err
}

func (db *SessionService) Create(s *wildcare.Session) error {
	_, err := db.Exec("INSERT INTO sessions (token, user_id, expires_at) VALUES(?, ?, ?)",
		s.Token, s.UserID, s.Expires)
	return err
}

func (db *SessionService) Delete(us *wildcare.Session) error {
	return nil
}

func (db *SessionService) Find(token string) (*wildcare.Session, error) {
	s := &wildcare.Session{}
	err := db.QueryRow(`
		SELECT user_id
		FROM sessions where token=?`, token).Scan(&s.UserID)

	return s, err
}
