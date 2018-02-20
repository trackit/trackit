import { put, call, all } from 'redux-saga/effects';
import { getToken } from '../misc';
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

export function* editAccountSaga({ account }) {
  try {
    const token = yield getToken();
    yield call(API.AWS.Accounts.editAccount, account, token);
    yield all([
      put({ type: Constants.AWS_EDIT_ACCOUNT_SUCCESS }),
      put({ type: Constants.AWS_GET_ACCOUNTS })
    ]);
  } catch (error) {
    yield put({ type: Constants.AWS_EDIT_ACCOUNT_ERROR, error });
  }
}

export function* editAccountBillSaga({ accountID, bill }) {
  try {
    const token = yield getToken();
    yield call(API.AWS.Accounts.editAccountBill, accountID, bill, token);
    yield all([
      put({ type: Constants.AWS_EDIT_ACCOUNT_BILL_SUCCESS }),
      put({ type: Constants.AWS_GET_ACCOUNT_BILLS, accountID })
    ]);
  } catch (error) {
    yield put({ type: Constants.AWS_EDIT_ACCOUNT_BILL_ERROR, error });
  }
}

export function* deleteAccountSaga({ accountID }) {
  try {
    const token = yield getToken();
    yield call(API.AWS.Accounts.deleteAccount, accountID, token);
    yield all([
      put({ type: Constants.AWS_DELETE_ACCOUNT_SUCCESS }),
      put({ type: Constants.AWS_GET_ACCOUNTS })
    ]);
  } catch (error) {
    yield put({ type: Constants.AWS_DELETE_ACCOUNT_ERROR, error });
  }
}

export function* deleteAccountBillSaga({ accountID, bill }) {
  try {
    const token = yield getToken();
    yield call(API.AWS.Accounts.deleteAccountBill, accountID, bill, token);
    yield all([
      put({ type: Constants.AWS_DELETE_ACCOUNT_BILL_SUCCESS }),
      put({ type: Constants.AWS_GET_ACCOUNTS })
    ]);
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
