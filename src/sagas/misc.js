import { select } from 'redux-saga/effects';
import Validations from '../common/forms/AWSAccountForm';

const getAccountIDFromRole = Validations.getAccountIDFromRole;

const getTokenFromState = (state) => (state.auth.token);

export const getToken = () => {
  return select(getTokenFromState);
};

const getAWSAccountsFromState = (state) => (state.aws.accounts.selection.map((account) => (getAccountIDFromRole(account.roleArn))));

export const getAWSAccounts = () => {
  return select(getAWSAccountsFromState);
};
