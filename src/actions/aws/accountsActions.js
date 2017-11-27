import Constants from '../../constants';

export default {
	getAccounts: () => ({
		type: Constants.AWS_GET_ACCOUNTS
	}),
  newAccount: (account) => ({
    type: Constants.AWS_NEW_ACCOUNT,
    account
  }),
  newAccountBill: (accountID, bill) => ({
    type: Constants.AWS_NEW_ACCOUNT_BILL,
    accountID,
    bill
  }),
  editAccount: (account) => ({
    type: Constants.AWS_EDIT_ACCOUNT,
    account
  }),
  editAccountBill: (accountID, bill) => ({
    type: Constants.AWS_EDIT_ACCOUNT_BILL,
    accountID,
    bill
  }),
  deleteAccount: (accountID) => ({
    type: Constants.AWS_DELETE_ACCOUNT,
    accountID
  }),
  deleteAccountBill: (accountID, bill) => ({
    type: Constants.AWS_DELETE_ACCOUNT_BILL,
    accountID,
    bill
  }),
	newExternal: () => ({
		type: Constants.AWS_NEW_EXTERNAL
	})
};
