package user

import "errors"

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Repo interface {
	Get(id int) (*User, error)
	Save(u *User) error
}

var (
	ErrUserNotFound = errors.New("User not found")
	ErrValidateUser = errors.New("Invalid user")
)
