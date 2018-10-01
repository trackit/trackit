import { call } from "../misc";

export const getEC2 = (token, accounts=undefined) => {
  let route = `/ec2`;
  if (accounts && accounts.length)
    route += `?accounts=${accounts.join(',')}`;
  return call(route, 'GET', null, token);
};

export const getEC2History = (token, date, accounts=undefined) => {
  let route = `/ec2/history?date=${date.format("YYYY-MM-DD")}`;
  if (accounts && accounts.length)
    route += `&accounts=${accounts.join(',')}`;
  return call(route, 'GET', null, token);
};

export const getRDS = (token, accounts=undefined) => {
  let route = `/rds`;
  if (accounts && accounts.length)
    route += `?accounts=${accounts.join(',')}`;
  return call(route, 'GET', null, token);
};

export const getRDSHistory = (token, date, accounts=undefined) => {
  let route = `/rds/history?date=${date.format("YYYY-MM-DD")}`;
  if (accounts && accounts.length)
    route += `&accounts=${accounts.join(',')}`;
  return call(route, 'GET', null, token);
};
