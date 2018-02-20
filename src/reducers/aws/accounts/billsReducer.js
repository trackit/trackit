import Constants from '../../../constants';

const defaultValue = {status: false};

export default (state=defaultValue, action) => {
  switch (action.type) {
    case Constants.AWS_GET_ACCOUNT_BILLS:
      return defaultValue;
    case Constants.AWS_GET_ACCOUNT_BILLS_SUCCESS:
      return {status: true, values: action.bills};
    case Constants.AWS_GET_ACCOUNT_BILLS_ERROR:
      return {status: true, error: action.error};
    case Constants.AWS_GET_ACCOUNT_BILLS_CLEAR:
      return {status: true, values: []};
    default:
      return state;
  }
};
