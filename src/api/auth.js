import { call } from './misc.js';

export const login = (email, password) => {
  return call('/user/login', 'POST', {email, password});
};

export const register = (email, password) => {
  return call('/user', 'POST', {email, password});
};
