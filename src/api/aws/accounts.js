import { call } from './../misc.js';

export const getAccounts = (token) => {
  return call('/aws', 'GET', null, token);
};

export const getAccountBills = (accountID, token) => {
  return call(`/aws/billrepository?aa=${accountID}`, 'GET', null, token);
};

export const newAccount = (account, token) => {
  return call('/aws', 'POST', account, token);
};

export const editAccount = (account, token) => {
//  return call(`/aws?aa=${account.id}`, 'PATCH', null, token);
  return { success: true, data: {} };
};

export const deleteAccount = (accountID, token) => {
//  return call(`/aws?aa=${accountID}`, 'DELETE', null, token);
  return { success: true, data: {} };
};

export const newAccountBill = (accountID, bill, token) => {
  return call(`/aws/billrepository?aa=${accountID}`, 'POST', bill, token);
};

export const editAccountBill = (accountID, bill, token) => {
//  return call(`/aws/billrepository?aa=${accountID}&br=${bill.id}`, 'PATCH', null, token);
  return { success: true, data: {} };
};

export const deleteAccountBill = (accountID, bill, token) => {
//  return call(`/aws/billrepository?aa=${accountID}&br=${bill.id}`, 'DELETE', null, token);
  return { success: true, data: {} };
};

export const newExternal = (token) => {
  return call('/aws/next', 'GET', null, token);
};
