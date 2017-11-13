import { takeLatest } from 'redux-saga/effects';
import { getAwsPricingSaga } from './pricingSaga';
import { getAccountsSaga, newAccountSaga, deleteAccountSaga, newExternalSaga } from './accountsSaga';
import { getS3DataSaga } from './s3Saga';

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

export function* watchDeleteAccount() {
  yield takeLatest(Constants.AWS_DELETE_ACCOUNT, deleteAccountSaga);
}

export function* watchGetAwsS3Data() {
  yield takeLatest(Constants.AWS_GET_S3_DATA, getS3DataSaga);
}
