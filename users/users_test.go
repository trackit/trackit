package users

// These tests are intended to be run against an empty database with the schema
// already in place.

import (
	"testing"
)

func TestNonExistingUserByIdFailure(t *testing.T) {
	min, max := uint(10000), uint(10100)
	for i := min; i < max; i++ {
		user, err := GetUserWithId(i)
		if user != nil {
			t.Errorf("User should be nil, instead points to %v.", *user)
		}
		if err == nil {
			t.Error("Error should not be nil, instead is nil.")
		} else if err.Error() != ErrorUserNotFound {
			t.Errorf("Error should be \"%s\", instead is \"%s\".", ErrorUserNotFound, err.Error())
		}
	}
}

func TestNonExistingUserByEmailFailure(t *testing.T) {
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
		user, err := GetUserWithEmail(email)
		if user != nil {
			t.Errorf("User should be nil, instead points to %v.", *user)
		}
		if err == nil {
			t.Error("Error should not be nil, instead is nil.")
		} else if err.Error() != ErrorUserNotFound {
			t.Errorf("Error should be \"%s\", instead is \"%s\".", ErrorUserNotFound, err.Error())
		}
	}
}

func TestCreateUpdateDelete(t *testing.T) {
	t.Run("CreateWithPasswordSuccess", testCreateWithPasswordSuccess)
	t.Run("GetAndPasswordSuccess", testGetAndPasswordSuccess)
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
	for _, c := range createWithPasswordCases {
		u := c.user
		if err := u.CreateWithPassword(c.password); err != nil {
			t.Errorf("Creating <%s>: error should be nil, instead is \"%s\".", u.Email, err.Error())
		} else {
			if u.Id <= 0 {
				t.Errorf("Creating <%s>: user ID should be >0, instead is %d.", u.Email, u.Id)
			}
		}
	}
}

func testGetAndPasswordSuccess(t *testing.T) {
	for _, c := range createWithPasswordCases {
		u, err := GetUserWithEmail(c.user.Email)
		if err != nil {
			t.Errorf("Getting <%s>: error should be nil, instead is \"%s\".", c.user.Email, err.Error())
		} else {
			m, err := u.PasswordMatches(c.password)
			if err != nil {
				t.Errorf("Checking good password for <%s>: error should be nil, instead is \"%s\".", u.Email, err.Error())
			} else if m == false {
				t.Error("Checking good password for <%s>: should match, doesn't.")
			}
			m, err = u.PasswordMatches(c.badPassword)
			if err != nil {
				t.Errorf("Checking bad password for <%s>: error should be nil, instead is \"%s\".", u.Email, err.Error())
			} else if m == false {
				t.Error("Checking bad password for <%s>: shouldn't match, does.")
			}
		}
	}
}
