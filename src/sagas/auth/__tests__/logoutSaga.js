import { put, all } from 'redux-saga/effects';
import logoutSaga from '../logoutSaga';
import Constants from '../../../constants';

const token = "42";
const mail = "mail";

describe("Logout Saga", () => {

  it("handles saga with available token", () => {

    let saga = logoutSaga();

    window.localStorage.setItem("userToken", token);
    window.localStorage.setItem("userMail", mail);

    expect(saga.next().value)
      .toEqual(all([
        put({ type: Constants.LOGOUT_REQUEST_SUCCESS }),
        put({ type: Constants.CLEAN_USER_TOKEN }),
        put({ type: Constants.CLEAN_USER_MAIL })
      ]));

    expect(saga.next().done).toBe(true);

  });

  it("handles saga with unavailable token", () => {

    let saga = logoutSaga();

    window.localStorage.setItem("userMail", mail);
    window.localStorage.removeItem("userToken");

    expect(saga.next().value)
      .toEqual(all([
        put({ type: Constants.LOGOUT_REQUEST_SUCCESS }),
        put({ type: Constants.CLEAN_USER_TOKEN }),
        put({ type: Constants.CLEAN_USER_MAIL })
      ]));

    expect(saga.next().done).toBe(true);

  });

  it("handles saga with unavailable mail", () => {

    let saga = logoutSaga();

    window.localStorage.setItem("userToken", token);
    window.localStorage.removeItem("userMail");

    expect(saga.next().value)
      .toEqual(all([
        put({ type: Constants.LOGOUT_REQUEST_SUCCESS }),
        put({ type: Constants.CLEAN_USER_TOKEN }),
        put({ type: Constants.CLEAN_USER_MAIL })
      ]));

    expect(saga.next().done).toBe(true);

  });

});
