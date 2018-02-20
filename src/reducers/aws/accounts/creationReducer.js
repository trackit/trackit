import Constants from '../../../constants';

const defaultValue = {status: false};

export default (state=defaultValue, action) => {
  switch (action.type) {
    case Constants.AWS_NEW_ACCOUNT:
      return defaultValue;
    case Constants.AWS_NEW_ACCOUNT_SUCCESS:
      return {status: true, value: action.account};
    case Constants.AWS_NEW_ACCOUNT_ERROR:
      return {status: true, error: action.error};
    case Constants.AWS_NEW_ACCOUNT_CLEAR:
      return {status: true, value: null};
    default:
      return state;
  }
};
