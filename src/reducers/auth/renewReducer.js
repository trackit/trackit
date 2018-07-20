import Constants from '../../constants';

export default (state=null, action) => {
  switch (action.type) {
    case Constants.RENEW_PASSWORD_SUCCESS:
      return { status: true, value: true };
    case Constants.RENEW_PASSWORD_ERROR:
      return { status: true, error: action.error };
    case Constants.RENEW_PASSWORD_LOADING:
      return {};
    case Constants.RENEW_PASSWORD_REQUEST:
    case Constants.RENEW_PASSWORD_CLEAR:
      return null;
    default:
      return state;
  }
};
