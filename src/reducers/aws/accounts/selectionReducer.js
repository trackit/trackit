import Constants from '../../../constants';

export default (state=[], action) => {
  switch (action.type) {
    case Constants.AWS_SELECT_ACCOUNT:
      let accounts = state.filter((item) => (item.id !== action.account.id));
      if (accounts.length === state.length)
        accounts.push(action.account);
      return accounts;
    case Constants.AWS_CLEAR_ACCOUNT_SELECTION:
      return [];
    default:
      return state;
  }
};
