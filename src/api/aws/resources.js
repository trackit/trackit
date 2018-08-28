import { call } from "../misc";

export const getEC2 = (token, accounts=undefined) => {
  let route = `/ec2`;
  if (accounts && accounts.length)
    route += `?accounts=${accounts.join(',')}`;
  return call(route, 'GET', null, token);
};

export const getRDS = (token, accounts=undefined) => {
  let route = `/rds`;
  if (accounts && accounts.length)
    route += `?accounts=${accounts.join(',')}`;
  return call(route, 'GET', null, token);
};
