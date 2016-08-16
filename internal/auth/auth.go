package auth

import (
	"github.com/larissavoigt/wildcare"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService struct {
}

func (a *AuthenticationService) HashPassword(u *wildcare.User, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hash)

	return nil
}

func (a *AuthenticationService) AuthenticateUser(u *wildcare.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
