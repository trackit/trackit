package users

import (
	"errors"
	"time"
)

const (
	ErrorNotImplemented = "Not implemented"
	ErrorUserNotFound   = "User not found"
)

type User struct {
	Id      uint
	Email   string
	Created time.Time
}

func (u *User) CreateWithPassword(password string) error {
	return errors.New(ErrorNotImplemented)
}

func (u User) Delete() error {
	return errors.New(ErrorNotImplemented)
}

func (u User) UpdatePassword(password string) error {
	return errors.New(ErrorNotImplemented)
}

func (u User) PasswordMatches(password string) (bool, error) {
	return false, errors.New(ErrorNotImplemented)
}

func GetUserWithId(id uint) (*User, error) {
	return nil, errors.New(ErrorNotImplemented)
}

func GetUserWithEmail(email string) (*User, error) {
	return nil, errors.New(ErrorNotImplemented)
}
