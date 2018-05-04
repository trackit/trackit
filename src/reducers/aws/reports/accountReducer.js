import Constants from '../../../constants';

export default (state='', action) => {
  switch (action.type) {
    case Constants.AWS_GET_ACCOUNTS_SUCCESS:
      let selectedId = '';
      if (action.accounts && action.accounts.length > 0) {
        selectedId = action.accounts[0].id.toString();
      }
      return selectedId
    case Constants.AWS_REPORTS_ACCOUNT_SELECTION:
      return action.accountId;
    default:
      return state;
  }
};
