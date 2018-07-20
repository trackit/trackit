import Constants from '../../constants';

export default (state=null, action) => {
  switch (action.type) {
    case Constants.RECOVER_PASSWORD_SUCCESS:
      return { status: true, value: true };
    case Constants.RECOVER_PASSWORD_ERROR:
      return { status: false, error: action.error };
    case Constants.RECOVER_PASSWORD_LOADING:
      return {};
    case Constants.RECOVER_PASSWORD_REQUEST:
    case Constants.RECOVER_PASSWORD_CLEAR:
      return null;
    default:
      return state;
  }
};
