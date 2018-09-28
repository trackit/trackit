import { call } from './../misc.js';

export const getAccounts = (token) => {
  return call('/aws?detailed', 'GET', null, token);
};

export const getAccountBills = (accountID, token) => {
  return call(`/aws/billrepository?account-id=${accountID}`, 'GET', null, token);
};

export const getAccountBillsStatus = (token) => {
  return call(`/aws/billrepositoryupdates`, 'GET', null, token);
};

export const newAccount = (account, token) => {
  return call('/aws', 'POST', account, token);
};

export const editAccount = (account, token) => {
  return call(`/aws?account-id=${account.id}`, 'PATCH', account, token);
};

export const deleteAccount = (accountID, token) => {
  return call(`/aws?account-id=${accountID}`, 'DELETE', null, token);
};

export const newAccountBill = (accountID, bill, token) => {
  return call(`/aws/billrepository?account-id=${accountID}`, 'POST', bill, token);
};

export const editAccountBill = (accountID, bill, token) => {
  return call(`/aws/billrepository?account-id=${accountID}&br=${bill.id}`, 'PATCH', bill, token);
};

export const deleteAccountBill = (accountID, billID, token) => {
  return call(`/aws/billrepository?account-id=${accountID}&br=${billID}`, 'DELETE', null, token);
};

export const newExternal = (token) => {
  return call('/aws/next', 'GET', null, token);
};

export const getAccountViewer = (accountID, token) => {
  return call(`/user/share?account-id=${accountID}`, 'GET', null, token);
};

export const addAccountViewer = (accountID, email, permissionLevel, token) => {
  return call(`/user/share?account-id=${accountID}`, 'POST', {email, permissionLevel}, token);
};

export const editAccountViewer = (sharedID, permissionLevel, token) => {
  return call(`/user/share?share-id=${sharedID}`, 'PATCH', {permissionLevel}, token);
};

export const deleteAccountViewer = (sharedID, token) => {
  return call(`/user/share?share-id=${sharedID}`, 'DELETE', null, token);
};
