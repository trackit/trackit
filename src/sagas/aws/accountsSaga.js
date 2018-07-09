import { put, call, all } from 'redux-saga/effects';
import { getSelectedAccounts, getToken } from '../misc';
import { setSelectedAccounts, getSelectedAccounts as getSelectedAccountsLS } from '../../common/localStorage';
import API from '../../api';
import Constants from '../../constants';

export function* getAccountsSaga() {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Accounts.getAccounts, token);
    if (res.success && res.hasOwnProperty("data"))
      yield put({ type: Constants.AWS_GET_ACCOUNTS_SUCCESS, accounts: res.data });
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_GET_ACCOUNTS_ERROR, error });
  }
}

export function* getAccountBillsSaga({ accountID }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Accounts.getAccountBills, accountID, token);
    if (res.success && res.hasOwnProperty("data"))
      yield put({ type: Constants.AWS_GET_ACCOUNT_BILLS_SUCCESS, bills: res.data });
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_GET_ACCOUNT_BILLS_ERROR, error });
  }
}

export function* newAccountSaga({ account }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Accounts.newAccount, account, token);
    if (res.success && res.hasOwnProperty("data"))
      yield all([
        put({ type: Constants.AWS_NEW_ACCOUNT_SUCCESS, account: res.data }),
        put({ type: Constants.AWS_NEW_EXTERNAL }),
        put({ type: Constants.AWS_GET_ACCOUNTS })
      ]);
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_NEW_ACCOUNT_ERROR, error });
  }
}

export function* newAccountBillSaga({ accountID, bill }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Accounts.newAccountBill, accountID, bill, token);
    if (res.success && res.hasOwnProperty("data"))
      yield all([
        put({ type: Constants.AWS_NEW_ACCOUNT_BILL_SUCCESS, bucket: res.data}),
        put({ type: Constants.AWS_GET_ACCOUNT_BILLS, accountID })
      ]);
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_NEW_ACCOUNT_BILL_ERROR, error });
  }
}

export function* getAccountBillStatusSaga() {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Accounts.getAccountBillsStatus, token);
    if (res.success && res.hasOwnProperty("data"))
      yield put({ type: Constants.AWS_GET_ACCOUNT_BILL_STATUS_SUCCESS, values: res.data });
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_GET_ACCOUNT_BILL_STATUS_ERROR, error });
  }
}

export function* editAccountSaga({ account }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Accounts.editAccount, account, token);
    if (res.success && res.hasOwnProperty("data"))
      yield all([
        put({ type: Constants.AWS_EDIT_ACCOUNT_SUCCESS }),
        put({ type: Constants.AWS_GET_ACCOUNTS })
      ]);
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_EDIT_ACCOUNT_ERROR, error });
  }
}

export function* editAccountBillSaga({ accountID, bill }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Accounts.editAccountBill, accountID, bill, token);
    if (res.success && res.hasOwnProperty("data"))
      yield all([
        put({ type: Constants.AWS_EDIT_ACCOUNT_BILL_SUCCESS }),
        put({ type: Constants.AWS_GET_ACCOUNT_BILLS, accountID })
      ]);
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_EDIT_ACCOUNT_BILL_ERROR, error });
  }
}

export function* deleteAccountSaga({ accountID }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Accounts.deleteAccount, accountID, token);
    if (res.success && res.hasOwnProperty("data"))
      yield all([
        put({ type: Constants.AWS_DELETE_ACCOUNT_SUCCESS }),
        put({ type: Constants.AWS_GET_ACCOUNTS })
      ]);
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_DELETE_ACCOUNT_ERROR, error });
  }
}

export function* deleteAccountBillSaga({ accountID, billID }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Accounts.deleteAccountBill, accountID, billID, token);
    if (res.success && res.hasOwnProperty("data"))
      yield all([
        put({ type: Constants.AWS_DELETE_ACCOUNT_BILL_SUCCESS }),
        put({ type: Constants.AWS_GET_ACCOUNT_BILLS, accountID })
      ]);
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_DELETE_ACCOUNT_BILL_ERROR, error });
  }
}

export function* newExternalSaga() {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Accounts.newExternal, token);
    if (res.success && res.hasOwnProperty("data") && res.data.hasOwnProperty("external") && res.data.hasOwnProperty("accountId"))
      yield put({ type: Constants.AWS_NEW_EXTERNAL_SUCCESS, external: res.data });
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_NEW_EXTERNAL_ERROR, error });
  }
}

export function* saveSelectedAccountSaga() {
  const data = yield getSelectedAccounts();
  setSelectedAccounts(data);
}

export function* loadSelectedAccountSaga() {
  try {
    const data = yield call(getSelectedAccountsLS);
    if (!data)
      throw Error("No selected accounts available");
    else if (Array.isArray(data))
      yield put({type: Constants.AWS_INSERT_SELECTED_ACCOUNTS, accounts: data});
    else
      throw Error("Invalid data for selected accounts");
  } catch (error) {
    yield put({ type: Constants.AWS_LOAD_SELECTED_ACCOUNTS_ERROR, error });
  }
}
