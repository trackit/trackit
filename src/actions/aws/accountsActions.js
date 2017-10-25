import Constants from '../../constants';

export default {
	getAccounts: () => ({
		type: Constants.AWS_GET_ACCOUNTS
	}),
	newAccount: (account) => ({
		type: Constants.AWS_NEW_ACCOUNT,
		account
	}),
	newExternal: () => ({
		type: Constants.AWS_NEW_EXTERNAL
	})
};
