import Constants from '../../../constants';

export default (state=[], action) => {
  switch (action.type) {
    case Constants.AWS_GET_ACCOUNT_BILLS_SUCCESS:
      return action.bills;
    case Constants.AWS_GET_ACCOUNT_BILLS_ERROR:
    case Constants.AWS_GET_ACCOUNT_BILLS_CLEAR:
      return [];
    default:
      return state;
  }
};
