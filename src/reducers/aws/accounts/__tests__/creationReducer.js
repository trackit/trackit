import CreationReducer from '../creationReducer';
import Constants from '../../../../constants';

const account = "account";

const defaultValue = {status: false};
const successValue = {status: true, value: account};
const errorValue = {status: true, error: Error()};
const cleared = {status: true, value: null};

describe("CreationReducer", () => {

  it("handles initial state", () => {
    expect(CreationReducer(undefined, {})).toEqual(defaultValue);
  });

  it("handles new account requested state", () => {
    expect(CreationReducer(null, { type: Constants.AWS_NEW_ACCOUNT })).toEqual(defaultValue);
  });

  it("handles new account success state", () => {
    expect(CreationReducer(defaultValue, { type: Constants.AWS_NEW_ACCOUNT_SUCCESS, account })).toEqual(successValue);
  });

  it("handles new account fail state", () => {
    expect(CreationReducer(successValue, { type: Constants.AWS_NEW_ACCOUNT_ERROR, error: errorValue.error })).toEqual(errorValue);
  });

  it("handles clear new account state", () => {
    expect(CreationReducer(successValue, { type: Constants.AWS_NEW_ACCOUNT_CLEAR })).toEqual(cleared);
  });

  it("handles wrong type state", () => {
    expect(CreationReducer(successValue, { type: "" })).toEqual(successValue);
  });

});
