import { put, call } from 'redux-saga/effects';
import registrationSaga from '../registrationSaga';
import API from '../../../api';
import Constants from '../../../constants';

describe("Registration Saga", () => {

  const credentials = { username: "username", password: "password" };
  const validResponse = { success: true, data: {} };
  const errorResponse = { success: true, data: { error: "error" } };
  const noResponse = { success: false };

  it("handles saga with valid data", () => {

    let saga = registrationSaga(credentials);

    expect(saga.next().value)
      .toEqual(put({ type: Constants.REGISTRATION_REQUEST_LOADING }));

    expect(saga.next().value)
      .toEqual(call(API.Auth.register, credentials.username, credentials.password));

    expect(saga.next(validResponse).value)
      .toEqual(put({ type: Constants.REGISTRATION_SUCCESS, payload: { status: true } }));

    expect(saga.next().done).toBe(true);

  });

  it("handles saga with no response", () => {

    let saga = registrationSaga(credentials);

    expect(saga.next().value)
      .toEqual(put({ type: Constants.REGISTRATION_REQUEST_LOADING }));

    expect(saga.next().value)
      .toEqual(call(API.Auth.register, credentials.username, credentials.password));

    expect(saga.next(noResponse).value)
      .toEqual(put({ type: Constants.REGISTRATION_ERROR, payload: { status: false, error: "Error: An error has occured" } }));

    expect(saga.next().done).toBe(true);

  });

  it("handles saga with error", () => {

    let saga = registrationSaga(credentials);

    expect(saga.next().value)
      .toEqual(put({ type: Constants.REGISTRATION_REQUEST_LOADING }));

    expect(saga.next().value)
      .toEqual(call(API.Auth.register, credentials.username, credentials.password));

    expect(saga.next(errorResponse).value)
      .toEqual(put({ type: Constants.REGISTRATION_ERROR, payload: { status: false, error: Error(errorResponse.data.error).toString() } }));

    expect(saga.next().done).toBe(true);

  });


});
