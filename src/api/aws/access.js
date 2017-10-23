import { call } from './../misc.js';

export const getAccess = (token) => {
  return call('/aws', 'GET', null, token);
};
