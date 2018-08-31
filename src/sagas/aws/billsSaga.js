import { put, call, all } from 'redux-saga/effects';
import { getToken } from '../misc';
import API from '../../api';
import Constants from '../../constants';

export function* getAccountBillsSaga({ accountID }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Accounts.getAccountBills, accountID, token);
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data"))
      yield put({ type: Constants.AWS_GET_ACCOUNT_BILLS_SUCCESS, bills: res.data });
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_GET_ACCOUNT_BILLS_ERROR, error });
  }
}

export function* newAccountBillSaga({ accountID, bill }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Accounts.newAccountBill, accountID, bill, token);
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data") && res.data.hasOwnProperty("id"))
      yield all([
        put({ type: Constants.AWS_NEW_ACCOUNT_BILL_SUCCESS, bucket: res.data}),
        put({ type: Constants.AWS_GET_ACCOUNT_BILLS, accountID }),
        put({ type: Constants.AWS_GET_ACCOUNTS }),
      ]);
    else if (res.success && res.hasOwnProperty("data") && res.data.hasOwnProperty("error"))
      throw Error(res.data.error);
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
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data"))
      yield put({ type: Constants.AWS_GET_ACCOUNT_BILL_STATUS_SUCCESS, values: res.data });
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_GET_ACCOUNT_BILL_STATUS_ERROR, error });
  }
}

export function* editAccountBillSaga({ accountID, bill }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Accounts.editAccountBill, accountID, bill, token);
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data") && res.data.hasOwnProperty("id"))
      yield all([
        put({ type: Constants.AWS_EDIT_ACCOUNT_BILL_SUCCESS }),
        put({ type: Constants.AWS_GET_ACCOUNT_BILLS, accountID })
      ]);
    else if (res.success && res.hasOwnProperty("data") && res.data.hasOwnProperty("error"))
      throw Error(res.data.error);
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_EDIT_ACCOUNT_BILL_ERROR, error });
  }
}

export function* deleteAccountBillSaga({ accountID, billID }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Accounts.deleteAccountBill, accountID, billID, token);
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data"))
      yield all([
        put({ type: Constants.AWS_DELETE_ACCOUNT_BILL_SUCCESS }),
        put({ type: Constants.AWS_GET_ACCOUNT_BILLS, accountID }),
        put({ type: Constants.AWS_GET_ACCOUNTS }),
      ]);
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_DELETE_ACCOUNT_BILL_ERROR, error });
  }
}
