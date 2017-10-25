import AllReducer from '../allReducer';
import Constants from '../../../../constants';

describe("AllReducer", () => {

  it("handle initial state", () => {
    expect(AllReducer(undefined, {})).toEqual([]);
  });

  it("handle get accounts success state", () => {
    const accounts = ["account1", "account2"];
    expect(AllReducer(null, { type: Constants.AWS_GET_ACCOUNTS_SUCCESS, accounts })).toEqual(accounts);
  });

  it("handle get accounts fail state", () => {
    const accounts = ["account1", "account2"];
    expect(AllReducer(accounts, { type: Constants.AWS_GET_ACCOUNTS_ERROR })).toEqual([]);
  });

  it("handle wrong type state", () => {
    const accounts = ["account1", "account2"];
    expect(AllReducer(accounts, { type: "" })).toEqual(accounts);
  });

});