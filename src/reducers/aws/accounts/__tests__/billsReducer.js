import BillsReducer from '../billsReducer';
import Constants from '../../../../constants';

const bills = ["bill1", "bill2"];

const defaultValue = {status: false};
const successValue = {status: true, values: bills};
const errorValue = {status: true, error: Error()};
const cleared = {status: true, values: []};

describe("BillsReducer", () => {

  it("handles initial state", () => {
    expect(BillsReducer(undefined, {})).toEqual(defaultValue);
  });

  it("handles get accounts requested state", () => {
    expect(BillsReducer(null, { type: Constants.AWS_GET_ACCOUNT_BILLS })).toEqual(defaultValue);
  });

  it("handles get accounts success state", () => {
    expect(BillsReducer(defaultValue, { type: Constants.AWS_GET_ACCOUNT_BILLS_SUCCESS, bills })).toEqual(successValue);
  });

  it("handles get accounts fail state", () => {
    expect(BillsReducer(successValue, { type: Constants.AWS_GET_ACCOUNT_BILLS_ERROR, error: errorValue.error })).toEqual(errorValue);
  });

  it("handles clear accounts state", () => {
    expect(BillsReducer(successValue, { type: Constants.AWS_GET_ACCOUNT_BILLS_CLEAR })).toEqual(cleared);
  });

  it("handles wrong type state", () => {
    expect(BillsReducer(successValue, { type: "" })).toEqual(successValue);
  });

});
