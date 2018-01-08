import Constants from '../../../constants';

export default (state=null, action) => {
  switch (action.type) {
    case Constants.AWS_NEW_ACCOUNT_SUCCESS:
      return action.account;
    case Constants.AWS_NEW_ACCOUNT_ERROR:
    case Constants.AWS_NEW_ACCOUNT_CLEAR:
    case Constants.AWS_NEW_ACCOUNT:
      return null;
    default:
      return state;
  }
};
