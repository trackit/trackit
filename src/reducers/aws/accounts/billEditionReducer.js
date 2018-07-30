import Constants from '../../../constants';

const defaultValue = {status: false};

export default (state=defaultValue, action) => {
  switch (action.type) {
    case Constants.AWS_EDIT_ACCOUNT_BILL:
      return defaultValue;
    case Constants.AWS_EDIT_ACCOUNT_BILL_CLEAR:
      return {status: true, value: null};
    case Constants.AWS_EDIT_ACCOUNT_BILL_SUCCESS:
      return {status: true, value: action.bucket};
    case Constants.AWS_EDIT_ACCOUNT_BILL_ERROR:
      return {status: true, error: action.error};
    default:
      return state;
  }
};
