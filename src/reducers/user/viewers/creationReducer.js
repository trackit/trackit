import Constants from '../../../constants';

const defaultValue = {status: false};

export default (state=defaultValue, action) => {
  switch (action.type) {
    case Constants.USER_NEW_VIEWER:
      return defaultValue;
    case Constants.USER_NEW_VIEWER_SUCCESS:
      return {status: true, value: action.viewer};
    case Constants.USER_NEW_VIEWER_ERROR:
      return {status: true, error: action.error};
    case Constants.USER_NEW_VIEWER_CLEAR:
      return {status: true, value: null};
    default:
      return state;
  }
};
