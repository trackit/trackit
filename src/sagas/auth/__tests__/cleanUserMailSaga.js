import { put } from 'redux-saga/effects';
import cleanUserMailSaga from '../cleanUserMailSaga';
import Constants from '../../../constants';

describe("Clean User Mail Saga", () => {

  it("handles saga", () => {

    let saga = cleanUserMailSaga();

    expect(saga.next().value)
      .toEqual(put({ type: Constants.CLEAN_USER_MAIL_SUCCESS}));

    expect(saga.next().done).toBe(true);

  });


});
