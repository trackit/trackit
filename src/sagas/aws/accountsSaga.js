import { put, call, all, select } from 'redux-saga/effects';
import API from '../../api';
import Constants from '../../constants';

const getToken = state => state.auth.token;

export function* getAccountsSaga() {
  try {
    const token = yield select(getToken);
    const res = yield call(API.AWS.Accounts.getAccounts, token);
    yield all([
      put({ type: Constants.AWS_GET_ACCOUNTS_SUCCESS, accounts: res.data }),
    ]);
  } catch (error) {
    yield put({ type: Constants.AWS_GET_ACCOUNTS_ERROR, error });
  }
}

export function* newAccountSaga({ account }) {
  try {
    const token = yield select(getToken);
    const res = yield call(API.AWS.Accounts.newAccount, account, token);
    yield all([
      put({ type: Constants.AWS_NEW_ACCOUNT_SUCCESS, account: res.data }),
      put({ type: Constants.AWS_GET_ACCOUNTS })
    ]);
  } catch (error) {
    yield put({ type: Constants.AWS_NEW_ACCOUNT_ERROR, error });
  }
}
