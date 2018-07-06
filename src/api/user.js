import { call } from './misc.js';

export const current = token => {
  return call('/user', 'GET', null, token);
};

export const Viewers = {
  list: token => call('/user/viewer', 'GET', null, token),
  create: (email, token) => call('/user/viewer', 'POST', { email }, token),
};
