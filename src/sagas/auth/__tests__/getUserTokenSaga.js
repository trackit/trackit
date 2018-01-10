import { put, call } from 'redux-saga/effects';
import getUserTokenSaga from '../getUserTokenSaga';
import { getToken } from "../../../common/localStorage";
import Constants from '../../../constants';

const token = "42";

describe("Get User Token Saga", () => {

  it("handles saga with available token", () => {

    let saga = getUserTokenSaga();

    window.localStorage.setItem("userToken", token);

    expect(saga.next().value)
      .toEqual(call(getToken));

    expect(saga.next(token).value)
      .toEqual(put({ type: Constants.GET_USER_TOKEN_SUCCESS , token}));

    expect(saga.next().done).toBe(true);

  });

  it("handles saga with unavailable token", () => {

    let saga = getUserTokenSaga();

    window.localStorage.removeItem("userToken");

    expect(saga.next().value)
      .toEqual(call(getToken));

    expect(saga.next().value)
      .toEqual(put({ type: Constants.GET_USER_TOKEN_ERROR, error: Error("No token available") }));

    expect(saga.next().done).toBe(true);

  });


});
