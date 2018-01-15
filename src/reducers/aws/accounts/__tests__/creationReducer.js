import CreationReducer from '../creationReducer';
import Constants from '../../../../constants';

describe("CreationReducer", () => {

  it("handles initial state", () => {
    expect(CreationReducer(undefined, {})).toEqual(null);
  });

  it("handles new account success state", () => {
    const account = "account";
    expect(CreationReducer(null, { type: Constants.AWS_NEW_ACCOUNT_SUCCESS, account })).toEqual(account);
  });

  it("handles new account fail state", () => {
    const account = "account";
    expect(CreationReducer(account, { type: Constants.AWS_NEW_ACCOUNT_ERROR })).toEqual(null);
  });

  it("handles clear new account state", () => {
    const account = "account";
    expect(CreationReducer(account, { type: Constants.AWS_NEW_ACCOUNT_CLEAR })).toEqual(null);
  });

  it("handles wrong type state", () => {
    const account = "account";
    expect(CreationReducer(account, { type: "" })).toEqual(account);
  });

});
