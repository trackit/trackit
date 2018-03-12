import AllReducer from '../allReducer';
import Constants from '../../../../constants';

const accounts = ["account1", "account2"];

const defaultValue = {status: false};
const successValue = {status: true, values: accounts};
const errorValue = {status: true, error: Error()};
const cleared = {status: true, values: []};

describe("AllReducer", () => {

  it("handles initial state", () => {
    expect(AllReducer(undefined, {})).toEqual(defaultValue);
  });

  it("handles get accounts requested state", () => {
    expect(AllReducer(null, { type: Constants.AWS_GET_ACCOUNTS })).toEqual(defaultValue);
  });

  it("handles get accounts success state", () => {
    expect(AllReducer(null, { type: Constants.AWS_GET_ACCOUNTS_SUCCESS, accounts })).toEqual(successValue);
  });

  it("handles get accounts fail state", () => {
    expect(AllReducer(defaultValue, { type: Constants.AWS_GET_ACCOUNTS_ERROR, error: errorValue.error })).toEqual(errorValue);
  });

  it("handles clear get accounts state", () => {
    expect(AllReducer(successValue, { type: Constants.AWS_GET_ACCOUNTS_CLEAR })).toEqual(cleared);
  });

  it("handles wrong type state", () => {
    expect(AllReducer(successValue, { type: "" })).toEqual(successValue);
  });

});
