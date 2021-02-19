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

package users

// These tests are intended to be run against an empty database with the schema
// already in place.

import (
	"context"
	"testing"

	"github.com/trackit/trackit/db"
)

func TestNonExistingUserByIdFailure(t *testing.T) {
	min, max := 10000, 10100
	for i := min; i < max; i++ {
		_, err := GetUserWithId(db.Db, i)
		if err == nil {
			t.Error("Error should not be nil, instead is nil.")
		} else if err != ErrUserNotFound {
			t.Errorf("Error should be \"%s\", instead is \"%s\".", ErrUserNotFound.Error(), err.Error())
		}
	}
}

func TestNonExistingUserByEmailFailure(t *testing.T) {
	ctx := context.Background()
	emails := []string{
		"",
		"bad format",
		"thisisbadlyformatted",
		"@nouser.com",
		"nodomain@",
		"notld@notld",
		"test.nexist@example.trackit.io",
		"toto.nexist@example.trackit.io",
		"lol.nexist@example.trackit.io",
	}
	for _, email := range emails {
		_, err := GetUserWithEmailAndOrigin(ctx, db.Db, email, "trackit")
		if err == nil {
			t.Error("Error should not be nil, instead is nil.")
		} else if err != ErrUserNotFound {
			t.Errorf("Error should be \"%s\", instead is \"%s\".", ErrUserNotFound.Error(), err.Error())
		}
	}
}

func TestCreateUpdateDelete(t *testing.T) {
	t.Run("CreateWithPasswordSuccess", testCreateWithPasswordSuccess)
	t.Run("GetAndPasswordSuccess", testGetAndPasswordSuccess)
	t.Run("GetAndPasswordFailure", testGetAndPasswordFailure)
}

type createWithPasswordCase struct {
	user        User
	password    string
	badPassword string
}

var createWithPasswordCases = [5]createWithPasswordCase{
	createWithPasswordCase{
		User{Email: "victor@example.trackit.io"},
		"testPassword",
		"badPassword",
	},
	createWithPasswordCase{
		User{Email: "victor_g@example.trackit.io"},
		"victor_gspassword",
		"",
	},
	createWithPasswordCase{
		User{Email: "lolwut@example.trackit.io"},
		"123456789",
		"321654987",
	},
	createWithPasswordCase{
		User{Email: "auietsrn@example.trackit.io"},
		"Pαßwòrd wìth Ûniķod€",
		"BäÐ §asswŏrd",
	},
	createWithPasswordCase{
		User{Email: "vdljqgh@example.trackit.io"},
		"", // Not the responsibility of this function to check for bad passwords.
		" ",
	},
}

func testCreateWithPasswordSuccess(t *testing.T) {
	ctx := context.Background()
	for _, c := range createWithPasswordCases {
		u, err := CreateUserWithPassword(ctx, db.Db, c.user.Email, c.password, "", "trackit")
		if err != nil {
			t.Errorf("Creating <%s>: error should be nil, instead is \"%s\".", c.user.Email, err.Error())
		} else {
			if u.Id <= 0 {
				t.Errorf("Creating <%s>: user ID should be >0, instead is %d.", c.user.Email, u.Id)
			}
		}
	}
}

func testGetAndPasswordSuccess(t *testing.T) {
	ctx := context.Background()
	for _, c := range createWithPasswordCases {
		u, err := GetUserFromOriginWithEmailAndPassword(ctx, db.Db, c.user.Email, c.password, "trackit")
		if err != nil {
			t.Errorf("Getting <%s>: error should be nil, instead is \"%s\".", c.user.Email, err.Error())
		}
		if u.Email != c.user.Email {
			t.Errorf("Getting <%s>: email should be <%s>, instead is <%s>.", c.user.Email, c.user.Email, u.Email)
		}
	}
}

func testGetAndPasswordFailure(t *testing.T) {
	ctx := context.Background()
	for _, c := range createWithPasswordCases {
		_, err := GetUserFromOriginWithEmailAndPassword(ctx, db.Db, c.user.Email, c.badPassword, "")
		if err == nil {
			t.Errorf("Getting <%s>: error should not be nil, instead is nil.", c.user.Email, err.Error())
		}
	}
}
