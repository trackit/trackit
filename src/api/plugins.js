import { call } from "./misc";

export const getData = (token, accounts=undefined) => {
  let route = `/plugins/results`;
  if (accounts && accounts.length)
    route += `?accounts=${accounts.join(',')}`;
  return call(route, 'GET', null, token);
};
