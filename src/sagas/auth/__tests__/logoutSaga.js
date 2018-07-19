import { put, all } from 'redux-saga/effects';
import logoutSaga from '../logoutSaga';
import Constants from '../../../constants';

const token = "42";
const mail = "mail";
const selectedAccounts = "1,2,3";

describe("Logout Saga", () => {

  it("handles saga with available token", () => {

    let saga = logoutSaga();

    window.localStorage.setItem("userToken", token);
    window.localStorage.setItem("userMail", mail);
    window.localStorage.setItem("selectedAccounts", selectedAccounts);

    expect(saga.next().value)
      .toEqual(all([
        put({ type: Constants.LOGOUT_REQUEST_SUCCESS }),
        put({ type: Constants.CLEAN_USER_TOKEN }),
        put({ type: Constants.CLEAN_USER_MAIL }),
        put({ type: Constants.CLEAN_USER_SELECTED_ACCOUNTS })
      ]));

    expect(saga.next().done).toBe(true);

  });

  it("handles saga with unavailable token", () => {

    let saga = logoutSaga();

    window.localStorage.setItem("userMail", mail);
    window.localStorage.removeItem("userToken");
    window.localStorage.setItem("selectedAccounts", selectedAccounts);

    expect(saga.next().value)
      .toEqual(all([
        put({ type: Constants.LOGOUT_REQUEST_SUCCESS }),
        put({ type: Constants.CLEAN_USER_TOKEN }),
        put({ type: Constants.CLEAN_USER_MAIL }),
        put({ type: Constants.CLEAN_USER_SELECTED_ACCOUNTS })
      ]));

    expect(saga.next().done).toBe(true);

  });

  it("handles saga with unavailable mail", () => {

    let saga = logoutSaga();

    window.localStorage.setItem("userToken", token);
    window.localStorage.removeItem("userMail");
    window.localStorage.setItem("selectedAccounts", selectedAccounts);

    expect(saga.next().value)
      .toEqual(all([
        put({ type: Constants.LOGOUT_REQUEST_SUCCESS }),
        put({ type: Constants.CLEAN_USER_TOKEN }),
        put({ type: Constants.CLEAN_USER_MAIL }),
        put({ type: Constants.CLEAN_USER_SELECTED_ACCOUNTS })
      ]));

    expect(saga.next().done).toBe(true);

  });

  it("handles saga with unavailable selected accounts", () => {

    let saga = logoutSaga();

    window.localStorage.setItem("userToken", token);
    window.localStorage.setItem("userMail", mail);
    window.localStorage.removeItem("selectedAccounts");

    expect(saga.next().value)
      .toEqual(all([
        put({ type: Constants.LOGOUT_REQUEST_SUCCESS }),
        put({ type: Constants.CLEAN_USER_TOKEN }),
        put({ type: Constants.CLEAN_USER_MAIL }),
        put({ type: Constants.CLEAN_USER_SELECTED_ACCOUNTS })
      ]));

    expect(saga.next().done).toBe(true);

  });

});
