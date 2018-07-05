import { call } from './misc.js';

export const login = (email, password) => {
  return call('/user/login', 'POST', {email, password});
};

export const register = (email, password) => {
  return call('/user', 'POST', {email, password});
};

export const recoverPassword = (email) => {
  return call('/user/recover', 'POST', {email});
};

export const renewPassword = (email, password, token) => {
  return call('/user/renew', 'POST', {email, password, token});
};
