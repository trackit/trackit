import { takeLatest } from 'redux-saga/effects';
import { getAwsPricingSaga } from './pricingSaga';
import { getAccountsSaga, newAccountSaga, newExternalSaga } from './accountsSaga';

import Constants from '../../constants';

export function* watchGetAwsPricing() {
  yield takeLatest(Constants.AWS_GET_PRICING, getAwsPricingSaga);
}

export function* watchGetAccounts() {
  yield takeLatest(Constants.AWS_GET_ACCOUNTS, getAccountsSaga);
}

export function* watchNewAccount() {
  yield takeLatest(Constants.AWS_NEW_ACCOUNT, newAccountSaga);
}

export function* watchNewExternal() {
  yield takeLatest(Constants.AWS_NEW_EXTERNAL, newExternalSaga);
}