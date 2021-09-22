//   Copyright 2017 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

// Package users implements the creation, usage and management of TrackIt user accounts
package users

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/models"
)

var (
	ErrNotImplemented = errors.New("Not implemented")
	ErrUserNotFound   = errors.New("User not found")
	ErrUserExists     = errors.New("User already exists")
	ErrFailedCreating = errors.New("Failed to create user")
)

// User is a user of the platform. It is different from models.User which is
// the database representation of a User.
type User struct {
	Id                     int    `json:"id"`
	Email                  string `json:"email"`
	NextExternal           string `json:"-"`
	ParentId               *int   `json:"parentId,omitempty"`
	AwsCustomerEntitlement bool   `json:"aws_customer_entitlement"`
}

// CreateUserWithPassword creates a user with an email and a password. A nil
// error indicates a success.
func CreateUserWithPassword(ctx context.Context, db models.XODB, email string, password string, customerIdentifier string, origin string) (User, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbUser := models.User{
		Email:                  email,
		AwsCustomerIdentifier:  customerIdentifier,
		AwsCustomerEntitlement: false,
		Created:                time.Now(),
		AccountType:            origin,
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
	return UserFromDbUser(dbUser), err
}

// CreateTagbotUser creates a tagbot user row associated with a trackit user
func CreateTagbotUser(ctx context.Context, db models.XODB, userId int, awsCustomerIdentifier string) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbUser := models.TagbotUser{
		UserID:                 userId,
		AwsCustomerIdentifier:  awsCustomerIdentifier,
		AwsCustomerEntitlement: false,
	}
	err := dbUser.Insert(db)
	if err != nil {
		logger.Error("Failed to create tagbot user.", err.Error())
	}
	return err
}

// CreateUserWithParent creates a viewer user with an email and a parent. A nil
// error indicates a success.
func CreateUserWithParent(ctx context.Context, db models.XODB, email string, parent User) (User, string, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbUser := models.User{
		Email:                  email,
		ParentUserID:           sql.NullInt64{int64(parent.Id), true},
		AwsCustomerEntitlement: true,
	}
	var user User
	var passRandom [12]byte
	_, err := rand.Read(passRandom[:])
	if err != nil {
		logger.Error("Failed to produce secure random password.", err.Error())
		return user, "", err
	}
	var passHuman [16]byte
	base64.StdEncoding.Encode(passHuman[:], passRandom[:])
	auth, err := getPasswordHash(string(passHuman[:]))
	if err != nil {
		logger.Error("Failed to create password hash.", err.Error())
		return user, "", err
	}
	dbUser.Auth = auth
	err = dbUser.Insert(db)
	if err != nil {
		logger.Error("Failed to insert viewer user in database.", err.Error())
		return user, "", err
	}
	user = UserFromDbUser(dbUser)
	return user, string(passHuman[:]), nil
}

func GetUsersByParent(ctx context.Context, db models.XODB, parent User) ([]User, error) {
	dbUsers, err := models.UsersByParentUserID(db, sql.NullInt64{int64(parent.Id), true})
	if err != nil {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to get viewer users by their parent.", err.Error())
		return nil, err
	}
	res := make([]User, len(dbUsers))
	for i := range dbUsers {
		res[i] = UserFromDbUser(*dbUsers[i])
	}
	return res, nil
}

func GetUserParent(ctx context.Context, db models.XODB, user User) (parent User, err error) {
	if user.ParentId == nil {
		jsonlog.LoggerFromContextOrDefault(ctx).Error("Tried to get parent of orphan user.", user)
		err = errors.New("Failed getting parent of orphan user.")
		return
	}
	parent, err = GetUserWithId(db, *user.ParentId)
	return
}

// UpdateUserWithPassword updates a user with an email and a password. A nil
// error indicates a success.
func UpdateUserWithPassword(ctx context.Context, tx *sql.Tx, dbUser *models.User, email string, password string) (User, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbUser.Email = email
	auth, err := getPasswordHash(password)
	if err != nil {
		logger.Error("Failed to create password hash.", err.Error())
	} else {
		dbUser.Auth = auth
		err = dbUser.Update(tx)
		if err != nil {
			logger.Error("Failed to update user.", err.Error())
		}
	}
	return UserFromDbUser(*dbUser), err
}

func (u User) UpdateNextExternal(ctx context.Context, db models.XODB) error {
	dbUser, err := models.UserByID(db, u.Id)
	if err == nil {
		if u.NextExternal == "" {
			dbUser.NextExternal.Valid = false
		} else {
			dbUser.NextExternal.Valid = true
			dbUser.NextExternal.String = u.NextExternal
		}
		return dbUser.Update(db)
	} else {
		return err
	}
}

// Delete deletes the user. A nil error indicates a success.
func (u User) Delete() error {
	return ErrNotImplemented
}

// UpdatePassword updates a user's password. A nil error indicates a success.
func (u User) UpdatePassword(db models.XODB, password string) error {
	dbUser, err := models.UserByID(db, u.Id)
	if err != nil {
		return err
	}
	auth, err := getPasswordHash(password)
	if err != nil {
		return err
	}
	dbUser.Auth = auth
	return dbUser.Update(db)
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
		user := UserFromDbUser(*dbUser)
		return user, nil
	}
}

// GetUserWithEmailAndOrigin retrieves the user with the given unique pair (Email, AccountType). A nil error
// indicates a success.
func GetUserWithEmailAndOrigin(ctx context.Context, db models.XODB, email string, accountType string) (User, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbUser, err := models.UserByEmailAccountType(db, email, accountType)
	if err == sql.ErrNoRows {
		return User{}, ErrUserNotFound
	} else if err != nil {
		logger.Error("Error getting user from database.", err.Error())
		return User{}, err
	} else {
		return UserFromDbUser(*dbUser), nil
	}
}

// GetUserFromOriginWithEmailAndPassword retrieves the user with the given unique Email
// and stored hash matching the given password. A nil eror indicates a success.
func GetUserFromOriginWithEmailAndPassword(ctx context.Context, db models.XODB, email string, password string, origin string) (User, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbUser, err := models.UserByEmailAccountType(db, email, origin)
	if err == sql.ErrNoRows {
		return User{}, ErrUserNotFound
	} else if err != nil {
		logger.Error("Error getting user from database.", err.Error())
		return User{}, err
	} else {
		if dbUser.AccountType != origin {
			return User{}, ErrUserNotFound
		}
		err = passwordMatchesHash(password, dbUser.Auth)
		return UserFromDbUser(*dbUser), err
	}
}

// UserFromDbUser builds a users.User from a models.User.
func UserFromDbUser(dbUser models.User) User {
	u := User{
		Id:                     dbUser.ID,
		Email:                  dbUser.Email,
		AwsCustomerEntitlement: dbUser.AwsCustomerEntitlement,
	}
	if dbUser.NextExternal.Valid {
		u.NextExternal = dbUser.NextExternal.String
	}
	if dbUser.ParentUserID.Valid {
		parentId := int(dbUser.ParentUserID.Int64)
		u.ParentId = &parentId
	}
	return u
}
