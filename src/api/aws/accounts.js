import { call } from './../misc.js';

export const getAccounts = (token) => {
  return call('/aws', 'GET', null, token);
};

export const newAccount = (account, token) => {
  return call('/aws', 'POST', account, token);
};
