import { put, call } from 'redux-saga/effects';
import Constants from '../../../constants';
import { getToken} from "../../misc";
import { getReportsSaga, clearReportsSaga, downloadReportSaga } from '../reportsSaga';
import API from '../../../api';

describe("Reports Saga", () => {

  const token = "421";
  const accountId = '42';

  describe("Get reports", () => {

    const validReponse = {
      success: true,
      data: ['myreports/myfile.xlsx',]
    }

    const error = Error('Unable to retrieve the list of reports');

    const invalidReponse = {
      success: false,
      error
    }

    it("handles saga with valid data", () => {

      let saga = getReportsSaga({ accountId });

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Reports.getReports, token, accountId));

      expect(saga.next(validReponse).value)
        .toEqual(put({ type: Constants.AWS_GET_REPORTS_SUCCESS, reports: validReponse.data, account: accountId }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = getReportsSaga({ accountId });

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Reports.getReports, token, accountId));

      expect(saga.next(invalidReponse).value)
        .toEqual(put({ type: Constants.AWS_GET_REPORTS_ERROR, error, account: accountId }));

      expect(saga.next().done).toBe(true);

    });
  });

  describe("Clear report", () => {
    it("handles saga", () => {

      let saga = clearReportsSaga();
      expect(saga.next().value)
        .toEqual(put({ type: Constants.AWS_CLEAR_REPORT }));

      expect(saga.next().done).toBe(true);

    });
  });

  describe("Download report", () => {

    const reportType = 'myreports';
    const fileName = 'myfile.xlsx';

    const validReponse = {
      success: true,
      data: '42'
    }

    const error = Error('Failed to download report file');

    const invalidReponse = {
      success: false,
      error
    }

    it("handles saga with valid data", () => {

      let saga = downloadReportSaga({ accountId, reportType, fileName });

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Reports.getReport, token, accountId, reportType, fileName));

      expect(saga.next(validReponse).value)
        .toEqual(put({ type: Constants.AWS_DOWNLOAD_REPORT_SUCCESS, account: accountId, reportType, fileName }));

      saga.next();

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = downloadReportSaga({ accountId, reportType, fileName });

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.Reports.getReport, token, accountId, reportType, fileName));

      expect(saga.next(invalidReponse).value)
        .toEqual(put({ type: Constants.AWS_DOWNLOAD_REPORT_ERROR, error, account: accountId, reportType, fileName }));

      expect(saga.next().done).toBe(true);

    });
  });

});
