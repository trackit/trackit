import Constants from '../../../constants';

const defaultValue = {status: false};

export default (state=defaultValue, action) => {
  switch (action.type) {
    case Constants.AWS_GET_ACCOUNT_BILL_STATUS:
      return defaultValue;
    case Constants.AWS_GET_ACCOUNT_BILL_STATUS_SUCCESS:
      return {status: true, values: action.values};
    case Constants.AWS_GET_ACCOUNT_BILL_STATUS_ERROR:
      return {status: true, error: action.error};
    case Constants.AWS_GET_ACCOUNT_BILL_STATUS_CLEAR:
      return {status: true, values: []};
    default:
      return state;
  }
};
