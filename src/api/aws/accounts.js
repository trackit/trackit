import { call } from './../misc.js';

export const getAccounts = (token) => {
  return call('/aws', 'GET', null, token);
};

export const newAccount = (account, token) => {
  return call('/aws', 'POST', account, token);
};

export const newExternal = (token) => {
//  return call('/aws/external', 'GET', null, token);
  return { status: true, data: { external: "external_test" } };
};
