import DownloadReducer from '../downloadReducer';
import Constants from '../../../../constants';

describe("AccountReducer", () => {

  const err = 'error';
  const defaultValue = {failed: false};
  const errorValue = {failed: true, error: err};

  it("handles initial state", () => {
    expect(DownloadReducer(undefined, {})).toEqual(defaultValue);
  });

  it("handles download report requested state", () => {
    expect(DownloadReducer(defaultValue, { type: Constants.AWS_DOWNLOAD_REPORT_REQUESTED })).toEqual(defaultValue);
  });

  it("handles download report success state", () => {
    expect(DownloadReducer(defaultValue, { type: Constants.AWS_DOWNLOAD_REPORT_SUCCESS })).toEqual(defaultValue);
  });

  it("handles download report error state", () => {
    expect(DownloadReducer(defaultValue, { type: Constants.AWS_DOWNLOAD_REPORT_ERROR, error: err })).toEqual(errorValue);
  });
});
