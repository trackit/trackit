import { call } from './misc.js';

export const login = (email, password, awsToken) => {
  return call('/user/login', 'POST', {email, password, awsToken});
};

export const register = (email, password, awsToken) => {
  return call('/user', 'POST', {email, password, awsToken});
};

export const recoverPassword = (email) => {
  return call('/user/password/forgotten', 'POST', {email});
};

export const renewPassword = (id, password, token) => {
  return call('/user/password/reset', 'POST', {id, password, token});
};
