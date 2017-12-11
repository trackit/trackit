import { put, call, all } from 'redux-saga/effects';
import {
  getAccountsSaga, newAccountSaga, editAccountSaga, deleteAccountSaga,
  newExternalSaga,
  getAccountBillsSaga, newAccountBillSaga, editAccountBillSaga, deleteAccountBillSaga
} from '../accountsSaga';
import { getToken } from '../../misc';
import API from '../../../api';
import Constants from '../../../constants';

const token = "42";

describe("Accounts Saga", () => {

  describe("Get Accounts", () => {

    const accounts = ["account1", "account2"];
    const validResponse = { success: true, data: accounts };
    const invalidResponse = { success: true, accounts };
    const noResponse = { success: false };

    it("handles saga with valid data", () => {

      let saga = getAccountsSaga();

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.getAccounts, token));

      expect(saga.next(validResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_ACCOUNTS_SUCCESS, accounts }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = getAccountsSaga();

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.getAccounts, token));

      expect(saga.next(invalidResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_ACCOUNTS_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with no response", () => {

      let saga = getAccountsSaga();

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.getAccounts, token));

      expect(saga.next(noResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_ACCOUNTS_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

  });

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

  describe("New Account", () => {

    const account = {roleArn: "roleArn"};
    const validResponse = { success: true, data: account };
    const invalidResponse = { success: true, account };
    const noResponse = { success: false };

    it("handles saga with valid data", () => {

      let saga = newAccountSaga({account});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.newAccount, account, token));

      expect(saga.next(validResponse).value)
        .toEqual(all([
          put({ type: Constants.AWS_NEW_ACCOUNT_SUCCESS, account }),
          put({ type: Constants.AWS_NEW_EXTERNAL }),
          put({ type: Constants.AWS_GET_ACCOUNTS })
        ]));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = newAccountSaga({account});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.newAccount, account, token));

      expect(saga.next(invalidResponse).value)
        .toEqual(put({ type: Constants.AWS_NEW_ACCOUNT_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with no response", () => {

      let saga = newAccountSaga({account});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.newAccount, account, token));

      expect(saga.next(noResponse).value)
        .toEqual(put({ type: Constants.AWS_NEW_ACCOUNT_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

  });

  describe("Edit Account", () => {

    it("handles saga", () => {

      const account = {roleArn: "roleArn"};

      let saga = editAccountSaga({account});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.editAccount, account, token));

      expect(saga.next().value)
        .toEqual(all([
          put({ type: Constants.AWS_EDIT_ACCOUNT_SUCCESS }),
          put({ type: Constants.AWS_GET_ACCOUNTS })
        ]));

      expect(saga.next().done).toBe(true);

    });

  });

  describe("New External", () => {

    const external = "external";
    const validResponse = { success: true, data: { external } };
    const invalidResponse = { success: true, external };
    const noResponse = { success: false };

    it("handles saga with valid data", () => {

      let saga = newExternalSaga();

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.newExternal, token));

      expect(saga.next(validResponse).value)
        .toEqual(all([
          put({ type: Constants.AWS_NEW_EXTERNAL_SUCCESS, external })
        ]));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = newExternalSaga();

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.newExternal, token));

      expect(saga.next(invalidResponse).value)
        .toEqual(put({ type: Constants.AWS_NEW_EXTERNAL_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with no response", () => {

      let saga = newExternalSaga();

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.newExternal, token));

      expect(saga.next(noResponse).value)
        .toEqual(put({ type: Constants.AWS_NEW_EXTERNAL_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

  });

  describe("Delete Account", () => {

    it("handles saga", () => {

      const accountID = 42;

      let saga = deleteAccountSaga({accountID});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.deleteAccount, accountID, token));

      expect(saga.next().value)
        .toEqual(all([
          put({ type: Constants.AWS_DELETE_ACCOUNT_SUCCESS }),
          put({ type: Constants.AWS_GET_ACCOUNTS })
        ]));

      expect(saga.next().done).toBe(true);

    });

  });

});

describe("Account Bills Saga", () => {

  describe("New Account Bill", () => {

    it("handles saga", () => {

      const accountID = 42;
      const bill = {bucket: "test"};

      let saga = newAccountBillSaga({accountID, bill});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.newAccountBill, accountID, bill, token));

      expect(saga.next().value)
        .toEqual(all([
          put({ type: Constants.AWS_NEW_ACCOUNT_BILL_SUCCESS }),
          put({ type: Constants.AWS_GET_ACCOUNTS })
        ]));

      expect(saga.next().done).toBe(true);

    });

  });

  describe("Edit Account Bill", () => {

    it("handles saga", () => {

      const accountID = 42;
      const bill = {bucket: "test"};

      let saga = editAccountBillSaga({accountID, bill});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.editAccountBill, accountID, bill, token));

      expect(saga.next().value)
        .toEqual(all([
          put({ type: Constants.AWS_EDIT_ACCOUNT_BILL_SUCCESS }),
          put({ type: Constants.AWS_GET_ACCOUNT_BILLS, accountID })
        ]));

      expect(saga.next().done).toBe(true);

    });

  });

  describe("Delete Account Bill", () => {

    it("handles saga", () => {

      const accountID = 42;
      const bill = {bucket: "test"};

      let saga = deleteAccountBillSaga({accountID, bill});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.deleteAccountBill, accountID, bill, token));

      expect(saga.next().value)
        .toEqual(all([
          put({ type: Constants.AWS_DELETE_ACCOUNT_BILL_SUCCESS }),
          put({ type: Constants.AWS_GET_ACCOUNTS })
        ]));

      expect(saga.next().done).toBe(true);

    });

  });

});
