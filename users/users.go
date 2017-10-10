package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit2/models"
)

var (
	ErrNotImplemented = errors.New("Not implemented")
	ErrUserNotFound   = errors.New("User not found")
	ErrUserExists     = errors.New("User already exists")
)

// User is a user of the platform. It is different from models.User which is
// the database representation of a User.
type User struct {
	Id    int
	Email string
}

// CreateUserWithPassword creates a user with an email and a password. A nil
// error indicates a success.
func CreateUserWithPassword(ctx context.Context, db models.XODB, email string, password string) (User, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbUser := models.User{
		Email: email,
	}
	auth, err := getPasswordHash(password)
	if err != nil {
		logger.Error("Failed to create password hash.", err.Error())
	} else {
		dbUser.Auth = auth
		err = dbUser.Insert(db)
		if err != nil {
			logger.Error("Failed to create user.", err.Error())
		}
	}
	return userFromDbUser(dbUser), err
}

// Delete deletes the user. A nil error indicates a success.
func (u User) Delete() error {
	return ErrNotImplemented
}

// UpdatePassword updates a user's password. A nil error indicates a success.
func (u User) UpdatePassword(password string) error {
	return ErrNotImplemented
}

// PasswordMatches tests whether a password matches a user's stored hash. A nil
// error indicates a match.
func (u User) PasswordMatches(password string) error {
	return ErrNotImplemented
}

// GetUserWithId retrieves the user with the given unique Id. A nil error
// indicates a success.
func GetUserWithId(db models.XODB, id int) (User, error) {
	dbUser, err := models.UserByID(db, id)
	if err == sql.ErrNoRows {
		user := User{}
		return user, ErrUserNotFound
	} else if err != nil {
		user := User{}
		return user, err
	} else {
		user := userFromDbUser(*dbUser)
		return user, nil
	}
}

// GetUserWithEmail retrieves the user with the given unique Email. A nil error
// indicates a success.
func GetUserWithEmail(ctx context.Context, db models.XODB, email string) (User, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbUser, err := models.UserByEmail(db, email)
	if err == sql.ErrNoRows {
		return User{}, ErrUserNotFound
	} else if err != nil {
		logger.Error("Error getting user from database.", err.Error())
		return User{}, err
	} else {
		return userFromDbUser(*dbUser), nil
	}
}

// GetUserWithEmailAndPassword retrieves the user with the given unique Email
// and stored hash matching the given password. A nil eror indicates a success.
func GetUserWithEmailAndPassword(ctx context.Context, db models.XODB, email string, password string) (User, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbUser, err := models.UserByEmail(db, email)
	if err == sql.ErrNoRows {
		return User{}, ErrUserNotFound
	} else if err != nil {
		logger.Error("Error getting user from database.", err.Error())
		return User{}, err
	} else {
		err = passwordMatchesHash(password, dbUser.Auth)
		return userFromDbUser(*dbUser), err
	}
}

// userFromDbUser builds a users.User from a models.User.
func userFromDbUser(dbUser models.User) User {
	return User{
		Id:    dbUser.ID,
		Email: dbUser.Email,
	}
}
