import { put, call, all } from 'redux-saga/effects';
import {
  getAccountsSaga, newAccountSaga, editAccountSaga, deleteAccountSaga,
  newExternalSaga,
  saveSelectedAccountSaga, loadSelectedAccountSaga
} from '../accountsSaga';
import {
  getSelectedAccounts as getSelectedAccountsLS
} from '../../../common/localStorage';
import {getSelectedAccounts, getToken} from '../../misc';
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
        .toEqual(call(API.AWS.Accounts.getAccountsStatus, token));

      expect(saga.next(validResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_ACCOUNTS_SUCCESS, accounts }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = getAccountsSaga();

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.getAccountsStatus, token));

      expect(saga.next(invalidResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_ACCOUNTS_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with no response", () => {

      let saga = getAccountsSaga();

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.getAccountsStatus, token));

      expect(saga.next(noResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_ACCOUNTS_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

  });

  describe("New Account", () => {

    const account = {roleArn: "roleArn", id: 42};
    const bill = { bucket: "test", prefix: ""};
    const validResponse = { success: true, data: account };
    const invalidResponse = { success: true, account };
    const noResponse = { success: false };

    it("handles saga with valid data", () => {

      let saga = newAccountSaga({account, bill});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.newAccount, account, token));

      expect(saga.next(validResponse).value)
        .toEqual(all([
          put({ type: Constants.AWS_NEW_ACCOUNT_SUCCESS, account }),
          put({ type: Constants.AWS_NEW_EXTERNAL }),
          put({ type: Constants.AWS_GET_ACCOUNTS }),
          put({ type : Constants.AWS_NEW_ACCOUNT_BILL, accountID: 42, bill})
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

    const account = {roleArn: "roleArn"};
    const validResponse = {success: true, data: account};
    const invalidResponse = {success: true, account};
    const noResponse = {success: false};

    it("handles saga with valid data", () => {

      let saga = editAccountSaga({account});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.editAccount, account, token));

      expect(saga.next(validResponse).value)
        .toEqual(all([
          put({type: Constants.AWS_EDIT_ACCOUNT_SUCCESS}),
          put({type: Constants.AWS_GET_ACCOUNTS})
        ]));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = editAccountSaga({account});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.editAccount, account, token));

      expect(saga.next(invalidResponse).value)
        .toEqual(put({type: Constants.AWS_EDIT_ACCOUNT_ERROR, error: Error("Error with request")}));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with no response", () => {

      let saga = editAccountSaga({account});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.editAccount, account, token));

      expect(saga.next(noResponse).value)
        .toEqual(put({type: Constants.AWS_EDIT_ACCOUNT_ERROR, error: Error("Error with request")}));

      expect(saga.next().done).toBe(true);

    });

  });

  describe("New External", () => {

    const external = {
      external: "external",
      accountId: "accountId"
    };
    const validResponse = { success: true, data: external };
    const invalidResponse = { success: true, external };
    const noResponse = { success: false };

    it("handles saga with valid data", () => {

      let saga = newExternalSaga();

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.newExternal, token));

      expect(saga.next(validResponse).value)
        .toEqual(put({ type: Constants.AWS_NEW_EXTERNAL_SUCCESS, external }));

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

    const accountID = 42;
    const validResponse = {success: true, data: {}};
    const invalidResponse = {success: true, accountID};
    const noResponse = {success: false};

    it("handles saga with valid data", () => {

      let saga = deleteAccountSaga({accountID});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.deleteAccount, accountID, token));

      expect(saga.next(validResponse).value)
        .toEqual(all([
          put({type: Constants.AWS_DELETE_ACCOUNT_SUCCESS}),
          put({type: Constants.AWS_GET_ACCOUNTS})
        ]));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = deleteAccountSaga({accountID});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.deleteAccount, accountID, token));

      expect(saga.next(invalidResponse).value)
        .toEqual(put({type: Constants.AWS_DELETE_ACCOUNT_ERROR, error: Error("Error with request")}));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with no response", () => {

      let saga = deleteAccountSaga({accountID});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Accounts.deleteAccount, accountID, token));

      expect(saga.next(noResponse).value)
        .toEqual(put({type: Constants.AWS_DELETE_ACCOUNT_ERROR, error: Error("Error with request")}));

      expect(saga.next().done).toBe(true);

    });

  });

});

describe("Selected Accounts Saga", () => {

  describe("Save Selected Accounts", () => {

    it("handles saga", () => {

      let saga = saveSelectedAccountSaga();

      expect(saga.next().value)
        .toEqual(getSelectedAccounts());

      expect(saga.next([]).done).toBe(true);

    });

  });

  describe("Load Selected Accounts", () => {

    const data = ["account1","account2"];

    const invalidData = {};

    it("handles saga with valid data", () => {

      let saga = loadSelectedAccountSaga();

      expect(saga.next().value)
        .toEqual(call(getSelectedAccountsLS));

      expect(saga.next(data).value)
        .toEqual(put({type: Constants.AWS_INSERT_SELECTED_ACCOUNTS, accounts: data}));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = loadSelectedAccountSaga();

      expect(saga.next().value)
        .toEqual(call(getSelectedAccountsLS));

      expect(saga.next(invalidData).value)
        .toEqual(put({ type: Constants.AWS_LOAD_SELECTED_ACCOUNTS_ERROR, error: Error("Invalid data for selected accounts") }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with no response", () => {

      let saga = loadSelectedAccountSaga();

      expect(saga.next().value)
        .toEqual(call(getSelectedAccountsLS));

      expect(saga.next(null).value)
        .toEqual(put({ type: Constants.AWS_LOAD_SELECTED_ACCOUNTS_ERROR, error: Error("No selected accounts available") }));

      expect(saga.next().done).toBe(true);

    });

  });

});
