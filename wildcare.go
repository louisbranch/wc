package wildcare

type User struct {
	ID       int64
	Name     string
	Email    string
	Password string
}

type UserService interface {
	User(id int64) (*User, error)
	Create(u *User) error
	Authenticate(email, password string) (*User, bool)
}
