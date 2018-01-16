import SelectionReducer from '../selectionReducer';
import Constants from '../../../../constants';

describe("SelectionReducer", () => {

  it("handles initial state", () => {
    expect(SelectionReducer(undefined, {})).toEqual([]);
  });

  it("handles select account state", () => {
    const account = "account1";
    expect(SelectionReducer([], { type: Constants.AWS_SELECT_ACCOUNT, account })).toEqual([account]);
    expect(SelectionReducer([account], { type: Constants.AWS_SELECT_ACCOUNT, account })).toEqual([]);
  });

  it("handles clear accounts state", () => {
    const accounts = ["account1", "account2"];
    expect(SelectionReducer(accounts, { type: Constants.AWS_CLEAR_ACCOUNT_SELECTION })).toEqual([]);
  });

  it("handles wrong type state", () => {
    const accounts = ["account1", "account2"];
    expect(SelectionReducer(accounts, { type: "" })).toEqual(accounts);
  });

});
