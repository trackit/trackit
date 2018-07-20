import { call } from './misc.js';

export const login = (email, password) => {
  return call('/user/login', 'POST', {email, password});
};

export const register = (email, password) => {
  return call('/user', 'POST', {email, password});
};

export const recoverPassword = (email) => {
  return call('/user/password/forgotten', 'POST', {email});
};

export const renewPassword = (id, password, token) => {
  return call('/user/password/reset', 'POST', {id, password, token});
};
