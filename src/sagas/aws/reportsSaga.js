import { put, call } from 'redux-saga/effects';
import API from '../../api';
import Constants from '../../constants';
import { getToken } from "../misc";

var FileSaver = require('file-saver');

export function* getReportsSaga({ accountId }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Reports.getReports, token, accountId);
    if (res.success && res.hasOwnProperty("data") && !res.data.hasOwnProperty("error"))
      yield put({ type: Constants.AWS_GET_REPORTS_SUCCESS, reports: res.data, account: accountId });
    else
      throw Error("Unable to retrieve the list of reports");
  } catch (error) {
    yield put({ type: Constants.AWS_GET_REPORTS_ERROR, error, account: accountId });
  }
}

export function* clearReportsSaga() {
  yield put({ type: Constants.AWS_CLEAR_REPORT });
}

export function* downloadReportSaga({ accountId, reportType, fileName }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Reports.getReport, token, accountId, reportType, fileName);
    if (res.success && res.hasOwnProperty("data") && !res.data.hasOwnProperty("error"))
      yield put({ type: Constants.AWS_DOWNLOAD_REPORT_SUCCESS, account: accountId, reportType, fileName });
    else
      throw Error("Failed to download report file");
    FileSaver.saveAs(res.data, fileName);
  } catch (error) {
    yield put({ type: Constants.AWS_DOWNLOAD_REPORT_ERROR, error, account: accountId, reportType, fileName });
  }
}
