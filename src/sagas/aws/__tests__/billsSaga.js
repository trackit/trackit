import { put, call, all } from 'redux-saga/effects';
import {getAccountBillsSaga, newAccountBillSaga, editAccountBillSaga, deleteAccountBillSaga,} from '../billsSaga';
import {getToken} from '../../misc';
import API from '../../../api';
import Constants from '../../../constants';

const token = "42";

describe("Account Bills Saga", () => {

  describe("Get Account Bills", () => {

    const accountID = 42;
    const bills = ["bill1", "bill2"];
    const validResponse = { success: true, data: bills };
    const invalidResponse = { success: true, bills };
    const noResponse = { success: false };

    it("handles saga with valid data", () => {

      let saga = getAccountBillsSaga({ accountID });

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.getAccountBills, accountID, token));

      expect(saga.next(validResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_ACCOUNT_BILLS_SUCCESS, bills }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = getAccountBillsSaga({ accountID });

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.getAccountBills, accountID, token));

      expect(saga.next(invalidResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_ACCOUNT_BILLS_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with no response", () => {

      let saga = getAccountBillsSaga({ accountID });

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.getAccountBills, accountID, token));

      expect(saga.next(noResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_ACCOUNT_BILLS_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

  });

  describe("New Account Bill", () => {

    const accountID = 42;
    const bill = {id: 21, bucket: "test"};
    const validResponse = { success: true, data: bill };
    const errorResponse = { success: true, data: {error: "Error"} };
    const invalidResponse = { success: true, bill };
    const noResponse = { success: false };

    it("handles saga", () => {

      let saga = newAccountBillSaga({accountID, bill});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.newAccountBill, accountID, bill, token));

      expect(saga.next(validResponse).value)
        .toEqual(all([
          put({ type: Constants.AWS_NEW_ACCOUNT_BILL_SUCCESS, bucket: bill }),
          put({ type: Constants.AWS_GET_ACCOUNT_BILLS, accountID }),
          put({ type: Constants.AWS_GET_ACCOUNTS }),
        ]));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with error in response data", () => {

      let saga = newAccountBillSaga({accountID, bill});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.newAccountBill, accountID, bill, token));

      expect(saga.next(errorResponse).value)
        .toEqual(put({ type: Constants.AWS_NEW_ACCOUNT_BILL_ERROR, error: Error(errorResponse.data.error) }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = newAccountBillSaga({accountID, bill});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.newAccountBill, accountID, bill, token));

      expect(saga.next(invalidResponse).value)
        .toEqual(put({ type: Constants.AWS_NEW_ACCOUNT_BILL_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with no response", () => {

      let saga = newAccountBillSaga({accountID, bill});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.newAccountBill, accountID, bill, token));

      expect(saga.next(noResponse).value)
        .toEqual(put({ type: Constants.AWS_NEW_ACCOUNT_BILL_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

  });

  describe("Edit Account Bill", () => {

    const accountID = 42;
    const bill = {id: 21, bucket: "test"};
    const validResponse = { success: true, data: bill };
    const errorResponse = { success: true, data: {error: "Error"} };
    const invalidResponse = { success: true, bill };
    const noResponse = { success: false };

    it("handles saga", () => {

      let saga = editAccountBillSaga({accountID, bill});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.editAccountBill, accountID, bill, token));

      expect(saga.next(validResponse).value)
        .toEqual(all([
          put({ type: Constants.AWS_EDIT_ACCOUNT_BILL_SUCCESS }),
          put({ type: Constants.AWS_GET_ACCOUNT_BILLS, accountID })
        ]));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with error in response data", () => {

      let saga = editAccountBillSaga({accountID, bill});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.editAccountBill, accountID, bill, token));

      expect(saga.next(errorResponse).value)
        .toEqual(put({ type: Constants.AWS_EDIT_ACCOUNT_BILL_ERROR, error: Error(errorResponse.data.error) }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = editAccountBillSaga({accountID, bill});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.editAccountBill, accountID, bill, token));

      expect(saga.next(invalidResponse).value)
        .toEqual(put({ type: Constants.AWS_EDIT_ACCOUNT_BILL_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with no response", () => {

      let saga = editAccountBillSaga({accountID, bill});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.editAccountBill, accountID, bill, token));

      expect(saga.next(noResponse).value)
        .toEqual(put({ type: Constants.AWS_EDIT_ACCOUNT_BILL_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

  });

  describe("Delete Account Bill", () => {

    const accountID = 42;
    const billID = 84;
    const validResponse = { success: true, data: {} };
    const invalidResponse = { success: true, billID };
    const noResponse = { success: false };

    it("handles saga", () => {

      let saga = deleteAccountBillSaga({accountID, billID});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.deleteAccountBill, accountID, billID, token));

      expect(saga.next(validResponse).value)
        .toEqual(all([
          put({ type: Constants.AWS_DELETE_ACCOUNT_BILL_SUCCESS }),
          put({ type: Constants.AWS_GET_ACCOUNT_BILLS, accountID }),
          put({ type: Constants.AWS_GET_ACCOUNTS }),
        ]));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = deleteAccountBillSaga({accountID, billID});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.deleteAccountBill, accountID, billID, token));

      expect(saga.next(invalidResponse).value)
        .toEqual(put({ type: Constants.AWS_DELETE_ACCOUNT_BILL_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with no response", () => {

      let saga = deleteAccountBillSaga({accountID, billID});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.deleteAccountBill, accountID, billID, token));

      expect(saga.next(noResponse).value)
        .toEqual(put({ type: Constants.AWS_DELETE_ACCOUNT_BILL_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

  });

});
