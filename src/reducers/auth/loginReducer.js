import Constants from '../../constants';

export default (state=null, action) => {
  switch (action.type) {
    case Constants.LOGIN_REQUEST_SUCCESS:
      return { status: true };
    case Constants.LOGIN_REQUEST_ERROR:
      return { status: false, error: action.error };
    case Constants.LOGIN_REQUEST_LOADING:
      return {};
    case Constants.LOGIN_REQUEST:
      return null;
    default:
      return state;
  }
};
