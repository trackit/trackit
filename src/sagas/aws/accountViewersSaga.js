import { put, call, all } from 'redux-saga/effects';
import { getToken } from '../misc';
import API from '../../api';
import Constants from '../../constants';

export function* getAccountViewers({ accountID }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Accounts.getAccountViewer, accountID, token);
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data"))
      yield put({ type: Constants.AWS_GET_ACCOUNT_VIEWERS_SUCCESS, accounts: res.data });
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_GET_ACCOUNT_VIEWERS_ERROR, error });
  }
}

export function* addAccountViewer({ email, accountID, permissionLevel }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Accounts.addAccountViewer, accountID, email, permissionLevel, token);
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data") && res.data.hasOwnProperty("error"))
      throw Error(res.data.error);
    else if (res.success && res.hasOwnProperty("data"))
      yield all([
        put({ type: Constants.AWS_ADD_ACCOUNT_VIEWER_SUCCESS, accounts: res.data }),
        put({ type: Constants.AWS_GET_ACCOUNT_VIEWERS, accountID }),
      ]);
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_ADD_ACCOUNT_VIEWER_ERROR, error });
  }
}

export function* editAccountViewer({ accountID, shareID, permissionLevel }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Accounts.editAccountViewer, shareID, permissionLevel, token);
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data") && res.data.hasOwnProperty("error"))
      throw Error(res.data.error);
    else if (res.success && res.hasOwnProperty("data"))
      yield all([
        put({ type: Constants.AWS_EDIT_ACCOUNT_VIEWER_SUCCESS, accounts: res.data }),
        put({ type: Constants.AWS_GET_ACCOUNT_VIEWERS, accountID }),
      ]);
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_EDIT_ACCOUNT_VIEWER_ERROR, error });
  }
}

export function* deleteAccountViewer({ accountID, shareID }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Accounts.deleteAccountViewer, shareID, token);
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data"))
      yield all([
        put({ type: Constants.AWS_DELETE_ACCOUNT_VIEWER_SUCCESS }),
        put({ type: Constants.AWS_GET_ACCOUNT_VIEWERS, accountID })
      ]);
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_DELETE_ACCOUNT_ERROR, error });
  }
}
