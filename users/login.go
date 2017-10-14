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
	"encoding/json"
	"errors"
	"net/http"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit2/db"
)

// loginRequestBody is the expected request body for the LogIn route handler.
type loginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// loginResponseBody is the response body in case LogIn succeeds.
type loginResponseBody struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

// LogIn handles users attempting to log in. It shall return a valid token the
// caller can then use to call other routes.
func LogIn(response http.ResponseWriter, request *http.Request) {
	var body loginRequestBody
	err := decodeRequestBody(request, &body)
	if err == nil && isLoginRequestBodyValid(body) {
		attemptLogInWithValidBody(response, request, body)
	} else {
		response.WriteHeader(400)
	}
}

// decodeRequestBody decodes a JSON request body and returns nil in case it
// could do so.
func decodeRequestBody(request *http.Request, structuredBody interface{}) error {
	return json.NewDecoder(request.Body).Decode(structuredBody)
}

// isLoginRequestBodyValid checks the validity of a log in request body.
func isLoginRequestBodyValid(body loginRequestBody) bool {
	return body.Email != "" && body.Password != ""
}

// attemptLogInWithValidBody tries to authenticate and log a user in using a
// validated login request.
func attemptLogInWithValidBody(response http.ResponseWriter, request *http.Request, body loginRequestBody) {
	logger := jsonlog.LoggerFromContextOrDefault(request.Context())
	user, err := GetUserWithEmailAndPassword(request.Context(), db.Db, body.Email, body.Password)
	if err == nil {
		logAuthenticatedUserIn(response, request, user)
	} else {
		logger.Warning("Authentication failure.", body)
		response.WriteHeader(403)
	}
}

// logAuthenticatedUserIn generates a token for a user that's already been
// authenticated.
func logAuthenticatedUserIn(response http.ResponseWriter, request *http.Request, user User) {
	logger := jsonlog.LoggerFromContextOrDefault(request.Context())
	token, err := generateToken(user)
	if err == nil {
		body := loginResponseBody{
			User:  user,
			Token: token,
		}
		response.WriteHeader(200)
		json.NewEncoder(response).Encode(body)
		logger.Info("User logged in.", user)
	} else {
		logger.Error("Failed to generate token.", err.Error())
		response.WriteHeader(500)
	}
}

// TestToken tests a token's validity. For a valid token, it returns the user
// the token belongs to.
func TestToken(response http.ResponseWriter, request *http.Request) {
	var err error
	logger := jsonlog.LoggerFromContextOrDefault(request.Context())
	if authorization := request.Header["Authorization"]; authorization != nil && len(authorization) == 1 {
		var user User
		tokenString := authorization[0]
		if user, err = testToken(tokenString); err == nil {
			response.WriteHeader(200)
			json.NewEncoder(response).Encode(user)
			return
		}
	} else {
		err = errors.New("Authorization header not found.")
	}
	logger.Warning("Failed testing token.", err.Error())
	response.WriteHeader(400)
}
