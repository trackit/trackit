import MailReducer from '../mailReducer';
import Constants from '../../../constants';

describe("MailReducer", () => {

  it("handles initial state", () => {
    expect(MailReducer(undefined, {})).toEqual(null);
  });

  it("handles get mail success state", () => {
    const mail = "mail";
    expect(MailReducer(null, { type: Constants.GET_USER_MAIL_SUCCESS, mail })).toEqual(mail);
  });

  it("handles get mail fail state", () => {
    expect(MailReducer("mail", { type: Constants.GET_USER_MAIL_ERROR })).toEqual(null);
  });

  it("handles clean mail success state", () => {
    expect(MailReducer("mail", { type: Constants.CLEAN_USER_MAIL_SUCCESS })).toEqual(null);
  });

  it("handles clean mail fail state", () => {
    expect(MailReducer("mail", { type: Constants.CLEAN_USER_MAIL_ERROR })).toEqual(null);
  });

  it("handles wrong type state", () => {
    expect(MailReducer("mail", { type: "" })).toEqual("mail");
  });

});
