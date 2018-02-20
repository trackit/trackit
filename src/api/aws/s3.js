import { call } from "../misc";

export const getData = (token, begin, end, accounts=undefined) => {
  let route = `/s3/costs?begin=${begin.format("YYYY-MM-DD")}&end=${end.format("YYYY-MM-DD")}`;
  if (accounts && accounts.length)
    route += `&aas=${accounts.join(',')}`;
  return call(route, 'GET', null, token);
};
