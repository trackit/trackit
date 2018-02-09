import RetrievedReducer from '../retrievedReducer';
import Constants from '../../../../constants';

describe("RetrievedReducer", () => {

  it("handles initial state", () => {
    expect(RetrievedReducer(undefined, {})).toEqual(false);
  });

  it("handles get accounts success state", () => {
    const accounts = ["account1", "account2"];
    expect(RetrievedReducer(undefined, { type: Constants.AWS_GET_ACCOUNTS_SUCCESS, accounts })).toEqual(true);
  });

});
