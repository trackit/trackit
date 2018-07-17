import ReportListReducer from '../reportListReducer';
import Constants from '../../../../constants';

describe("AccountReducer", () => {

  const err = 'error';
  const reports = ['myreport/myfile.xlsx',];
  const defaultValue = {status: false, values: []};
  const successValue = {status: true, values: reports};
  const errorValue = {status: true, error: err};

  it("handles initial state", () => {
    expect(ReportListReducer(undefined, {})).toEqual(defaultValue);
  });

  it("handles clear report state", () => {
    expect(ReportListReducer(defaultValue, { type: Constants.AWS_CLEAR_REPORT })).toEqual(defaultValue);
  });

  it("handles get reports success state", () => {
    expect(ReportListReducer(defaultValue, { type: Constants.AWS_GET_REPORTS_SUCCESS, reports })).toEqual(successValue);
  });

  it("handles download report error state", () => {
    expect(ReportListReducer(defaultValue, { type: Constants.AWS_GET_REPORTS_ERROR, error: err })).toEqual(errorValue);
  });
});
