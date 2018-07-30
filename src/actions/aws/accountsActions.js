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
  clearNewAccountBill: () => ({
    type: Constants.AWS_NEW_ACCOUNT_BILL_CLEAR
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
  clearEditAccountBill: () => ({
    type: Constants.AWS_EDIT_ACCOUNT_BILL_CLEAR
  }),
  deleteAccount: (accountID) => ({
    type: Constants.AWS_DELETE_ACCOUNT,
    accountID
  }),
  deleteAccountBill: (accountID, billID) => ({
    type: Constants.AWS_DELETE_ACCOUNT_BILL,
    accountID,
    billID
  }),
	newExternal: () => ({
		type: Constants.AWS_NEW_EXTERNAL
	}),
  selectAccount: (account) => ({
    type: Constants.AWS_SELECT_ACCOUNT,
    account
  }),
  clearAccountSelection: () => ({
    type: Constants.AWS_CLEAR_ACCOUNT_SELECTION
  }),
  getAccountBillsStatus: () => ({
    type: Constants.AWS_GET_ACCOUNT_BILL_STATUS,
  }),
  clearAccountBillsStatus: () => ({
    type: Constants.AWS_GET_ACCOUNT_BILL_STATUS_CLEAR
  })
};
