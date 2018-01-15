import Constants from '../../constants';

export default (state=null, action) => {
  switch (action.type) {
    case Constants.GET_USER_MAIL_SUCCESS:
      return action.mail;
    case Constants.CLEAN_USER_MAIL_SUCCESS:
    case Constants.CLEAN_USER_MAIL_ERROR:
    case Constants.GET_USER_MAIL_ERROR:
      return null;
    default:
      return state;
  }
};
