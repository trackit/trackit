import Constants from '../../../constants';
import Validation from '../../../common/forms/AWSAccountForm';

export default (state='', action) => {
  switch (action.type) {
    case Constants.AWS_GET_ACCOUNTS_SUCCESS:
      if (action.accounts && action.accounts.length > 0)
        return Validation.getAccountIDFromRole(action.accounts[0].roleArn);
      return '';
    case Constants.AWS_RESOURCES_ACCOUNT_SELECTION:
      return action.accountId;
    default:
      return state;
  }
};
