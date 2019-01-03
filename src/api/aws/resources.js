import { call } from "../misc";

export const getEC2 = (token, date, accounts=undefined) => {
  let route = `/ec2?date=${date.format("YYYY-MM-DD")}`;
  if (accounts && accounts.length)
    route += `&accounts=${accounts.join(',')}`;
  return call(route, 'GET', null, token);
};

export const getUnusedEC2 = (token, date, accounts=undefined) => {
  let route = `/ec2/unused?date=${date.format("YYYY-MM-DD")}`;
  if (accounts && accounts.length)
    route += `&accounts=${accounts.join(',')}`;
  return call(route, 'GET', null, token);
};

export const getRDS = (token, date, accounts=undefined) => {
  let route = `/rds?date=${date.format("YYYY-MM-DD")}`;
  if (accounts && accounts.length)
    route += `&accounts=${accounts.join(',')}`;
  return call(route, 'GET', null, token);
};

export const getUnusedRDS = (token, date, accounts=undefined) => {
  let route = `/rds/unused?date=${date.format("YYYY-MM-DD")}`;
  if (accounts && accounts.length)
    route += `&accounts=${accounts.join(',')}`;
  return call(route, 'GET', null, token);
};

export const getES = (token, date, accounts=undefined) => {
  let route = `/es?date=${date.format("YYYY-MM-DD")}`;
  if (accounts && accounts.length)
    route += `&accounts=${accounts.join(',')}`;
  return call(route, 'GET', null, token);
};

export const getUnusedES = (token, date, accounts=undefined) => {
  let route = `/es/unused?date=${date.format("YYYY-MM-DD")}`;
  if (accounts && accounts.length)
    route += `&accounts=${accounts.join(',')}`;
  return call(route, 'GET', null, token);
};

export const getElastiCache = (token, date, accounts=undefined) => {
  let route = `/elasticache?date=${date.format("YYYY-MM-DD")}`;
  if (accounts && accounts.length)
    route += `&accounts=${accounts.join(',')}`;
  return call(route, 'GET', null, token);
};

export const getUnusedElastiCache = (token, date, accounts=undefined) => {
  let route = `/elasticache/unused?date=${date.format("YYYY-MM-DD")}`;
  if (accounts && accounts.length)
    route += `&accounts=${accounts.join(',')}`;
  return call(route, 'GET', null, token);
};

export const getLambdas = (token, date, accounts=undefined) => {
  let route = `/lambda?date=${date.format("YYYY-MM-DD")}`;
  if (accounts && accounts.length)
    route += `&accounts=${accounts.join(',')}`;
  return call(route, 'GET', null, token);
};