import AccountReducer from '../accountReducer';
import Constants from '../../../../constants';

describe("AccountReducer", () => {

  const accounts = [
    {id: 42},
    {id: 420}
  ];

  const defaultValue = '';
  const defaultAccount = '42';
  const selectedAccountId = '420';

  it("handles initial state", () => {
    expect(AccountReducer(undefined, {})).toEqual(defaultValue);
  });

  it("handles empty get account requested state", () => {
    expect(AccountReducer(defaultValue, { type: Constants.AWS_GET_ACCOUNTS_SUCCESS, accounts: [] })).toEqual('');
  });

  it("handles get account requested state", () => {
    expect(AccountReducer(defaultValue, { type: Constants.AWS_GET_ACCOUNTS_SUCCESS, accounts })).toEqual(defaultAccount);
  });

  it("handles account selection requested state", () => {
    expect(AccountReducer(defaultValue, { type: Constants.AWS_REPORTS_ACCOUNT_SELECTION, accountId: selectedAccountId })).toEqual(selectedAccountId);
  });
});
