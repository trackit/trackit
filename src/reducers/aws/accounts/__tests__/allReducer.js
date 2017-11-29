import AllReducer from '../allReducer';
import Constants from '../../../../constants';

describe("AllReducer", () => {

  it("handles initial state", () => {
    expect(AllReducer(undefined, {})).toEqual([]);
  });

  it("handles get accounts success state", () => {
    const accounts = ["account1", "account2"];
    expect(AllReducer(null, { type: Constants.AWS_GET_ACCOUNTS_SUCCESS, accounts })).toEqual(accounts);
  });

  it("handles get accounts fail state", () => {
    const accounts = ["account1", "account2"];
    expect(AllReducer(accounts, { type: Constants.AWS_GET_ACCOUNTS_ERROR })).toEqual([]);
  });

  it("handles wrong type state", () => {
    const accounts = ["account1", "account2"];
    expect(AllReducer(accounts, { type: "" })).toEqual(accounts);
  });

});
