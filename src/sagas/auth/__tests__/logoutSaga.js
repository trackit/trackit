import { put, all } from 'redux-saga/effects';
import logoutSaga from '../logoutSaga';
import Constants from '../../../constants';

const token = "42";

describe("Logout Saga", () => {

  it("handless saga with available token", () => {

    let saga = logoutSaga();

    window.localStorage.setItem("userToken", token);

    expect(saga.next().value)
      .toEqual(all([
        put({ type: Constants.LOGOUT_REQUEST_SUCCESS }),
        put({ type: Constants.CLEAN_USER_TOKEN })
      ]));

    expect(saga.next().done).toBe(true);

    expect(window.localStorage.getItem("userToken")).toBe(null);

  });

  it("handless saga with unavailable token", () => {

    let saga = logoutSaga();

    window.localStorage.removeItem("userToken");

    expect(saga.next().value)
      .toEqual(put({ type: Constants.LOGOUT_REQUEST_ERROR, error: Error("No token available") }));

    expect(saga.next().done).toBe(true);

    expect(window.localStorage.getItem("userToken")).toBe(null);

  });

});
