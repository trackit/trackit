import { put, call } from 'redux-saga/effects';
import getUserMailSaga from '../getUserMailSaga';
import { getUserMail } from "../../../common/localStorage";
import Constants from '../../../constants';

const mail = "mail";

describe("Get User Mail Saga", () => {

  it("handles saga with available token", () => {

    let saga = getUserMailSaga();

    window.localStorage.setItem("userMail", mail);

    expect(saga.next().value)
      .toEqual(call(getUserMail));

    expect(saga.next(mail).value)
      .toEqual(put({ type: Constants.GET_USER_MAIL_SUCCESS , mail}));

    expect(saga.next().done).toBe(true);

  });

  it("handles saga with unavailable token", () => {

    let saga = getUserMailSaga();

    window.localStorage.removeItem("userMail");

    expect(saga.next().value)
      .toEqual(call(getUserMail));

    expect(saga.next().value)
      .toEqual(put({ type: Constants.GET_USER_MAIL_ERROR, error: Error("No user mail available") }));

    expect(saga.next().done).toBe(true);

  });


});
