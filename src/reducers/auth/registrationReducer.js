import Constants from '../../constants';

export default (state=null, action) => {
  switch (action.type) {
    case Constants.REGISTRATION_SUCCESS:
    case Constants.REGISTRATION_ERROR:
      return action.payload;
    case Constants.REGISTRATION_CLEAR:
    case Constants.REGISTRATION_REQUEST:
      return null;
    default:
      return state;
  }
};
