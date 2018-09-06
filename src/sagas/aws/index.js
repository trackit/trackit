import { takeEvery, takeLatest } from 'redux-saga/effects';
import * as AccountsSaga from './accountsSaga';
import { getCostsSaga, saveChartsSaga, loadChartsSaga, initChartsSaga } from "./costsSaga";
import { getS3DataSaga, saveS3DatesSaga, loadS3DatesSaga } from './s3Saga';
import { getReportsSaga, clearReportsSaga, downloadReportSaga } from './reportsSaga';
import { getEC2ReportSaga, getRDSReportSaga } from './resourcesSaga';
import { getMapCostsSaga } from './mapSaga';
import { getTagsKeysSaga, getTagsValuesSaga, initTagsChartsSaga, loadTagsChartsSaga, saveTagsChartsSaga } from './tagsSaga';
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
  yield takeLatest(Constants.AWS_DOWNLOAD_REPORT_REQUESTED, downloadReportSaga);
}

export function* watchGetEC2Report() {
  yield takeLatest(Constants.AWS_RESOURCES_GET_EC2, getEC2ReportSaga);
}

export function* watchGetRDSReport() {
  yield takeLatest(Constants.AWS_RESOURCES_GET_RDS, getRDSReportSaga);
}

export function* watchGetMapCosts() {
  yield takeLatest(Constants.AWS_MAP_GET_COSTS, getMapCostsSaga);
}

export function* watchGetTagsKeys() {
  yield takeEvery(Constants.AWS_TAGS_GET_KEYS, getTagsKeysSaga);
}

export function* watchGetTagsValues() {
  yield takeEvery(Constants.AWS_TAGS_GET_VALUES, getTagsValuesSaga);
}

export function* watchInitTagsCharts() {
  yield takeLatest(Constants.AWS_TAGS_INIT_CHARTS, initTagsChartsSaga);
}

export function* watchLoadTagsCharts() {
  yield takeLatest(Constants.AWS_TAGS_LOAD_CHARTS, loadTagsChartsSaga);
}

export function* watchSaveTagsCharts() {
  yield takeEvery(Constants.AWS_TAGS_ADD_CHART, saveTagsChartsSaga);
  yield takeEvery(Constants.AWS_TAGS_REMOVE_CHART, saveTagsChartsSaga);
  yield takeEvery(Constants.AWS_TAGS_SET_DATES, saveTagsChartsSaga);
  yield takeEvery(Constants.AWS_TAGS_RESET_DATES, saveTagsChartsSaga);
  yield takeEvery(Constants.AWS_TAGS_SET_INTERVAL, saveTagsChartsSaga);
  yield takeEvery(Constants.AWS_TAGS_CLEAR_INTERVAL, saveTagsChartsSaga);
}

export function* watchGetAccountBillStatus() {
  yield takeLatest(Constants.AWS_GET_ACCOUNT_BILL_STATUS, AccountsSaga.getAccountBillStatusSaga)
}
