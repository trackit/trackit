import { takeLatest } from 'redux-saga/effects';
import { getS3DataSaga } from './s3Saga';
import { getAccountsSaga, newAccountSaga, newExternalSaga } from './accountsSaga';

import Constants from '../../constants';

export function* watchGetAccounts() {
  yield takeLatest(Constants.AWS_GET_ACCOUNTS, getAccountsSaga);
}

export function* watchNewAccount() {
  yield takeLatest(Constants.AWS_NEW_ACCOUNT, newAccountSaga);
}

export function* watchNewExternal() {
  yield takeLatest(Constants.AWS_NEW_EXTERNAL, newExternalSaga);
}

export function* watchGetAwsS3Data() {
  yield takeLatest(Constants.AWS_GET_S3_DATA, getS3DataSaga);
}
