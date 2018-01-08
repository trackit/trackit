import Constants from '../../constants';

export default {
  getAccounts: () => ({
    type: Constants.AWS_GET_ACCOUNTS
  }),
  clearAccounts: () => ({
    type: Constants.AWS_GET_ACCOUNTS_CLEAR
  }),
  getAccountBills: (accountID) => ({
    type: Constants.AWS_GET_ACCOUNT_BILLS,
    accountID
  }),
  clearAccountBills: () => ({
    type: Constants.AWS_GET_ACCOUNT_BILLS_CLEAR
  }),
  newAccount: (account) => ({
    type: Constants.AWS_NEW_ACCOUNT,
    account
  }),
  clearNewAccount: () => ({
    type: Constants.AWS_NEW_ACCOUNT_CLEAR
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
