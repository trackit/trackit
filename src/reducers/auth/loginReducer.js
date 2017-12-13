import Constants from '../../constants';

export default (state=null, { type, error }) => {
  switch (type) {
    case Constants.LOGIN_REQUEST_SUCCESS:
      return { status: true };
    case Constants.LOGIN_REQUEST_ERROR:
      return { status: false,  error};
    case Constants.LOGIN_REQUEST:
      return null;
    default:
      return state;
  }
};
