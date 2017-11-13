import Constants from '../../constants';

export default {
	getAccounts: () => ({
		type: Constants.AWS_GET_ACCOUNTS
	}),
  newAccount: (account) => ({
    type: Constants.AWS_NEW_ACCOUNT,
    account
  }),
  deleteAccount: (accountID) => ({
    type: Constants.AWS_DELETE_ACCOUNT,
    accountID
  }),
	newExternal: () => ({
		type: Constants.AWS_NEW_EXTERNAL
	})
};
