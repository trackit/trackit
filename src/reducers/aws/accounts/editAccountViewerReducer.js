import Constants from '../../../constants';

const defaultValue = {status: true, value: null};

export default (state=defaultValue, action) => {
  switch (action.type) {
    case Constants.AWS_EDIT_ACCOUNT_VIEWER_CLEAR:
      return defaultValue;
    case Constants.AWS_EDIT_ACCOUNT_VIEWER:
      return {status: false};
    case Constants.AWS_EDIT_ACCOUNT_VIEWER_SUCCESS:
      return {status: true, values: action.accounts};
    case Constants.AWS_EDIT_ACCOUNT_VIEWER_ERROR:
      return {status: true, error: action.error};
    default:
      return state;
  }
};
