import { put, call, all } from 'redux-saga/effects';
import loginSaga from '../loginSaga';
import API from '../../../api';
import Constants from '../../../constants';

const token = "42";
const email = "mail";

describe("Login Saga", () => {

  const credentials = { username: "username", password: "password" };
  const validResponse = { success: true, data: { token, user: { email } } };
  const validErrorResponse = { success: true, data: { error: "error" }};
  const invalidResponse = { success: true, token };
  const noResponse = { success: false };

  it("handles saga with valid data", () => {

    let saga = loginSaga(credentials);

    expect(saga.next().value)
      .toEqual(put({ type: Constants.LOGIN_REQUEST_LOADING }));

    expect(saga.next().value)
      .toEqual(call(API.Auth.login, credentials.username, credentials.password));

    expect(saga.next(validResponse).value)
      .toEqual(all([
        put({ type: Constants.LOGIN_REQUEST_SUCCESS }),
        put({ type: Constants.GET_USER_TOKEN }),
        put({ type: Constants.GET_USER_MAIL })
      ]));

    expect(saga.next().done).toBe(true);

    expect(window.localStorage.getItem("userToken")).toBe(token);
    expect(window.localStorage.getItem("userMail")).toBe(email);

  });

  it("handles saga with valid data when login error", () => {

    let saga = loginSaga(credentials);

    expect(saga.next().value)
      .toEqual(put({ type: Constants.LOGIN_REQUEST_LOADING }));

    expect(saga.next().value)
      .toEqual(call(API.Auth.login, credentials.username, credentials.password));

    expect(saga.next(validErrorResponse).value)
      .toEqual(put({ type: Constants.LOGIN_REQUEST_ERROR, error: validErrorResponse.data.error }));

    expect(saga.next().done).toBe(true);

    expect(window.localStorage.getItem("userToken")).toBe(token);

  });

  it("handles saga with invalid data", () => {

    let saga = loginSaga(credentials);

    expect(saga.next().value)
      .toEqual(put({ type: Constants.LOGIN_REQUEST_LOADING }));

    expect(saga.next().value)
      .toEqual(call(API.Auth.login, credentials.username, credentials.password));

    expect(saga.next(invalidResponse).value)
      .toEqual(put({ type: Constants.LOGIN_REQUEST_ERROR, error: "Error with request" }));

    expect(saga.next().done).toBe(true);

  });

  it("handles saga with no response", () => {

    let saga = loginSaga(credentials);

    expect(saga.next().value)
      .toEqual(put({ type: Constants.LOGIN_REQUEST_LOADING }));

    expect(saga.next().value)
      .toEqual(call(API.Auth.login, credentials.username, credentials.password));

    expect(saga.next(noResponse).value)
      .toEqual(put({ type: Constants.LOGIN_REQUEST_ERROR, error: "Error with request" }));

    expect(saga.next().done).toBe(true);

  });


});
