import { takeLatest } from 'redux-saga/effects';
import { getAwsPricingSaga } from './pricingSaga';
import * as AccountsSaga from './accountsSaga';
import { getS3DataSaga } from './s3Saga';

import Constants from '../../constants';

export function* watchGetAwsPricing() {
  yield takeLatest(Constants.AWS_GET_PRICING, getAwsPricingSaga);
}

export function* watchGetAccounts() {
  yield takeLatest(Constants.AWS_GET_ACCOUNTS, AccountsSaga.getAccountsSaga);
}

export function* watchNewAccount() {
  yield takeLatest(Constants.AWS_NEW_ACCOUNT, AccountsSaga.newAccountSaga);
}

export function* watchNewAccountBill() {
  yield takeLatest(Constants.AWS_NEW_ACCOUNT_BILL, AccountsSaga.newAccountBillSaga);
}

export function* watchEditAccount() {
  yield takeLatest(Constants.AWS_EDIT_ACCOUNT, AccountsSaga.editAccountSaga);
}

export function* watchEditAccountBill() {
  yield takeLatest(Constants.AWS_EDIT_ACCOUNT_BILL, AccountsSaga.editAccountBillSaga);
}

export function* watchDeleteAccount() {
  yield takeLatest(Constants.AWS_DELETE_ACCOUNT, AccountsSaga.deleteAccountSaga);
}

export function* watchDeleteAccountBill() {
  yield takeLatest(Constants.AWS_DELETE_ACCOUNT_BILL, AccountsSaga.deleteAccountBillSaga);
}

export function* watchNewExternal() {
  yield takeLatest(Constants.AWS_NEW_EXTERNAL, AccountsSaga.newExternalSaga);
}

export function* watchGetAwsS3Data() {
  yield takeLatest(Constants.AWS_GET_S3_DATA, getS3DataSaga);
  yield takeLatest(Constants.AWS_SET_S3_VIEW_DATES, getS3DataSaga);
}
