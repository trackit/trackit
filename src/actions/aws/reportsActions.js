import Constants from '../../constants';

export default {
	selectAccount: (accountId) => ({
		type: Constants.AWS_REPORTS_ACCOUNT_SELECTION,
		accountId
	})
};
