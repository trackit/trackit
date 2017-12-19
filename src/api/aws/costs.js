import { call } from './../misc.js';

export const getCosts = (token, begin, end, filters, accounts=undefined) => {
  let route = `/costs?begin=${begin.format("YYYY-MM-DD")}&end=${end.format("YYYY-MM-DD")}&by=${filters.join(',')}`;
  if (accounts)
    route += `&accounts=${accounts.join(',')}`;
  return call(route, 'GET', null, token);
};
