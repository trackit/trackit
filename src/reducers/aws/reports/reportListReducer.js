import Constants from '../../../constants';

const defaultValue = {status: false, values: []};

export default (state=defaultValue, action) => {
  switch (action.type) {
    case Constants.AWS_CLEAR_REPORT:
      return defaultValue;
    case Constants.AWS_GET_REPORTS_SUCCESS:
      return {status: true, values: action.reports};
    case Constants.AWS_GET_REPORTS_ERROR:
      return {status: true, error: action.error};
    default:
      return state;
  }
};
