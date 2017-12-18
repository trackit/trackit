import BillsReducer from '../billsReducer';
import Constants from '../../../../constants';

describe("BillsReducer", () => {

  it("handles initial state", () => {
    expect(BillsReducer(undefined, {})).toEqual([]);
  });

  it("handles get accounts success state", () => {
    const bills = ["bill1", "bill2"];
    expect(BillsReducer(null, { type: Constants.AWS_GET_ACCOUNT_BILLS_SUCCESS, bills })).toEqual(bills);
  });

  it("handles get accounts fail state", () => {
    const bills = ["bill1", "bill2"];
    expect(BillsReducer(bills, { type: Constants.AWS_GET_ACCOUNT_BILLS_ERROR })).toEqual([]);
  });

  it("handles clear accounts state", () => {
    const bills = ["bill1", "bill2"];
    expect(BillsReducer(bills, { type: Constants.AWS_GET_ACCOUNT_BILLS_CLEAR })).toEqual([]);
  });

  it("handles wrong type state", () => {
    const bills = ["bill1", "bill2"];
    expect(BillsReducer(bills, { type: "" })).toEqual(bills);
  });

});
