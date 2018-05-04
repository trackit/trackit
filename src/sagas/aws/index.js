import { takeEvery, takeLatest } from 'redux-saga/effects';
import * as AccountsSaga from './accountsSaga';
import { getCostsSaga, saveChartsSaga, loadChartsSaga, initChartsSaga } from "./costsSaga";
import { getS3DataSaga, saveS3DatesSaga, loadS3DatesSaga } from './s3Saga';
import { getReportsSaga, clearReportsSaga, downloadReportSaga } from './reportsSaga';
import Constants from '../../constants';

export function* watchGetAccounts() {
  yield takeLatest(Constants.AWS_GET_ACCOUNTS, AccountsSaga.getAccountsSaga);
  yield takeLatest(Constants.AWS_GET_ACCOUNTS, clearReportsSaga);
}

export function* watchGetAccountBills() {
  yield takeLatest(Constants.AWS_GET_ACCOUNT_BILLS, AccountsSaga.getAccountBillsSaga);
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

export function* watchSaveSelectedAccounts() {
  yield takeEvery(Constants.AWS_SELECT_ACCOUNT, AccountsSaga.saveSelectedAccountSaga);
  yield takeEvery(Constants.AWS_CLEAR_ACCOUNT_SELECTION, AccountsSaga.saveSelectedAccountSaga);
}

export function* watchLoadSelectedAccounts() {
  yield takeLatest(Constants.AWS_LOAD_SELECTED_ACCOUNTS, AccountsSaga.loadSelectedAccountSaga);
}

export function* watchGetCosts() {
  yield takeEvery(Constants.AWS_GET_COSTS, getCostsSaga);
}

export function* watchSaveCharts() {
  yield takeEvery(Constants.AWS_ADD_CHART, saveChartsSaga);
  yield takeEvery(Constants.AWS_REMOVE_CHART, saveChartsSaga);
  yield takeEvery(Constants.AWS_SET_COSTS_DATES, saveChartsSaga);
  yield takeEvery(Constants.AWS_RESET_COSTS_DATES, saveChartsSaga);
  yield takeEvery(Constants.AWS_SET_COSTS_INTERVAL, saveChartsSaga);
  yield takeEvery(Constants.AWS_RESET_COSTS_INTERVAL, saveChartsSaga);
  yield takeEvery(Constants.AWS_SET_COSTS_FILTER, saveChartsSaga);
  yield takeEvery(Constants.AWS_RESET_COSTS_FILTER, saveChartsSaga);
}

export function* watchLoadCharts() {
  yield takeLatest(Constants.AWS_LOAD_CHARTS, loadChartsSaga);
}

export function* watchInitCharts() {
  yield takeLatest(Constants.AWS_INIT_CHARTS, initChartsSaga);
}

export function* watchGetAwsS3Data() {
  yield takeLatest(Constants.AWS_GET_S3_DATA, getS3DataSaga);
}

export function* watchSaveS3Dates() {
  yield takeEvery(Constants.AWS_SET_S3_DATES, saveS3DatesSaga);
  yield takeEvery(Constants.AWS_CLEAR_S3_DATES, saveS3DatesSaga);
}

export function* watchLoadS3Data() {
  yield takeLatest(Constants.AWS_LOAD_S3_DATES, loadS3DatesSaga);
}

export function* watchGetReports() {
  yield takeLatest(Constants.AWS_GET_REPORTS_REQUESTED, getReportsSaga);
}

export function* watchSelectReports() {
  yield takeLatest(Constants.AWS_REPORTS_ACCOUNT_SELECTION, clearReportsSaga);
}

export function* watchDownloadReport() {
  yield takeLatest(Constants.AWS_DOWNLOAD_REPORT_REQUESTED, downloadReportSaga)
}
