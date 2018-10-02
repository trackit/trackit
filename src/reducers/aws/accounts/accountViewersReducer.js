import Constants from '../../../constants';

const defaultValue = {status: true, value: null};

export default (state=defaultValue, action) => {
  switch (action.type) {
    case Constants.AWS_GET_ACCOUNT_VIEWERS_CLEAR:
      return defaultValue;
    case Constants.AWS_GET_ACCOUNT_VIEWERS:
      return {status: false};
    case Constants.AWS_GET_ACCOUNT_VIEWERS_SUCCESS:
      return {status: true, values: action.accounts};
    case Constants.AWS_GET_ACCOUNT_VIEWERS_ERROR:
      return {status: true, error: action.error};
    default:
      return state;
  }
};
