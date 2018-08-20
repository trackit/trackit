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

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/trackit/trackit-server/config"
	"golang.org/x/crypto/bcrypt"
)

const (
	bCryptCost = 12
)

var (
	jwtIssuer                  = config.AuthIssuer
	jwtSecret                  = []byte(config.AuthSecret)
	ErrInvalidClaims           = errors.New("claims are invalid")
	ErrCannotReadToken         = errors.New("failed to read token")
	ErrMissingToken            = errors.New("missing or duplicate token")
	ErrFailedToValidateToken   = errors.New("failed to validate token")
	ErrMarketplaceInvalidToken = errors.New("failed to validate marketplace token")
)

// getPasswordHash generates a hash string for a given password.
func getPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bCryptCost)
	return string(hash), err
}

// passwordMatchesHash checks whether a password matches a hash.
func passwordMatchesHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// jwtClaims represents the JWT claims used by this software, as a structure.
type jwtClaims struct {
	Issuer    string `json:"iss"`
	NotBefore int64  `json:"nbf"`
	Expires   int64  `json:"exp"`
	Subject   int    `json:"sub"`
	User      User   `json:"usr"`
	jwt.StandardClaims
}

// generateToken generates a valid JWT token for a given user.
func generateToken(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims{
		Issuer:    jwtIssuer,
		NotBefore: time.Now().Add(-1 * time.Hour).Unix(),
		Expires:   time.Now().Add(60 * 24 * time.Hour).Unix(),
		Subject:   user.Id,
		User:      user,
	})
	return token.SignedString([]byte(jwtSecret))
}

// getTokenSigningKey is used by jwt-go to check whether a token is acceptable
// before verifying it.
func getTokenSigningKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v.", token.Header["alg"])
	} else {
		return jwtSecret, nil
	}
}

// areClaimsValid checks whether the claims of a JWT token make it currently
// valid.
func areClaimsValid(claims jwtClaims) bool {
	now := time.Now().Unix()
	return claims.Issuer == jwtIssuer && claims.NotBefore <= now && now < claims.Expires
}

// testToken checks whether a JWT token is valid and retrieves the owning User
// if it is.
func testToken(tx *sql.Tx, tokenString string) (User, error) {
	var user User
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, getTokenSigningKey)
	if err == nil {
		if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
			if areClaimsValid(*claims) {
				userId := claims.Subject
				user, err = GetUserWithId(tx, userId)
				if user.AwsCustomerEntitlement {
					err = ErrMarketplaceInvalidToken
				}
			} else {
				err = ErrInvalidClaims
			}
		} else {
			err = ErrCannotReadToken
		}
	} else {
		err = ErrCannotReadToken
	}
	return user, err
}
