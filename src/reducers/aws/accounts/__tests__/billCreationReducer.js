import BillCreationReducer from '../billCreationReducer';
import Constants from '../../../../constants';

const bill = "bill";

const defaultValue = {status: false};
const successValue = {status: true, value: bill};
const errorValue = {status: true, error: Error()};
const cleared = {status: true, value: null};

describe("BillCreationReducer", () => {

  it("handles initial state", () => {
    expect(BillCreationReducer(undefined, {})).toEqual(defaultValue);
  });

  it("handles new bill requested state", () => {
    expect(BillCreationReducer(null, { type: Constants.AWS_NEW_ACCOUNT_BILL })).toEqual(defaultValue);
  });

  it("handles new bill success state", () => {
    expect(BillCreationReducer(defaultValue, { type: Constants.AWS_NEW_ACCOUNT_BILL_SUCCESS, bucket: bill })).toEqual(successValue);
  });

  it("handles new bill fail state", () => {
    expect(BillCreationReducer(successValue, { type: Constants.AWS_NEW_ACCOUNT_BILL_ERROR, error: errorValue.error })).toEqual(errorValue);
  });

  it("handles new bill clear state", () => {
    expect(BillCreationReducer(successValue, { type: Constants.AWS_NEW_ACCOUNT_BILL_CLEAR })).toEqual(cleared);
  });

  it("handles wrong type state", () => {
    expect(BillCreationReducer(successValue, { type: "" })).toEqual(successValue);
  });

});
