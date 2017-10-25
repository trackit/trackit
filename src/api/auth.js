import { call } from './misc.js';

export const login = (email, password) => {
  return call('/login', 'POST', {email, password});
};
