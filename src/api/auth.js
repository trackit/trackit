import { call } from './misc.js';

export const login = (email, password) => {
  return call('/user/login', 'POST', {email, password});
};
