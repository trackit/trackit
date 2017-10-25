import Constants from '../../../constants';

export default (state=[], action) => {
  switch (action.type) {
    case Constants.AWS_GET_ACCOUNTS_SUCCESS:
      return action.accounts;
    default:
      return state;
  }
};
