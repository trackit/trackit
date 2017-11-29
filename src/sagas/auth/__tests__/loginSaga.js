import { put, call, all } from 'redux-saga/effects';
import loginSaga from '../loginSaga';
import API from '../../../api';
import Constants from '../../../constants';

const token = "42";

describe("Login Saga", () => {

  const credentials = { username: "username", password: "password" };
  const validResponse = { success: true, data: { token } };
  const invalidResponse = { success: true, token };
  const noResponse = { success: false };

  it("handles saga with valid data", () => {

    let saga = loginSaga(credentials);

    expect(saga.next().value)
      .toEqual(call(API.Auth.login, credentials.username, credentials.password));

    expect(saga.next(validResponse).value)
      .toEqual(all([
        put({ type: Constants.LOGIN_REQUEST_SUCCESS }),
        put({ type: Constants.GET_USER_TOKEN })
      ]));

    expect(saga.next().done).toBe(true);

    expect(window.localStorage.getItem("userToken")).toBe(token);

  });

  it("handles saga with invalid data", () => {

    let saga = loginSaga(credentials);

    expect(saga.next().value)
      .toEqual(call(API.Auth.login, credentials.username, credentials.password));

    expect(saga.next(invalidResponse).value)
      .toEqual(put({ type: Constants.LOGIN_REQUEST_ERROR, error: Error("Error with request") }));

    expect(saga.next().done).toBe(true);

  });

  it("handles saga with no response", () => {

    let saga = loginSaga(credentials);

    expect(saga.next().value)
      .toEqual(call(API.Auth.login, credentials.username, credentials.password));

    expect(saga.next(noResponse).value)
      .toEqual(put({ type: Constants.LOGIN_REQUEST_ERROR, error: Error("Error with request") }));

    expect(saga.next().done).toBe(true);

  });


});
