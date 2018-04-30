import Constants from '../../../constants';

const defaultValue = {failed: false};

export default (state=defaultValue, action) => {
  switch (action.type) {
    case Constants.AWS_DOWNLOAD_REPORT_REQUESTED:
    case Constants.AWS_DOWNLOAD_REPORT_SUCCESS:
      return defaultValue;
    case Constants.AWS_DOWNLOAD_REPORT_ERROR:
      return {failed: true, error: action.error}
    default:
      return state;
  }
};
