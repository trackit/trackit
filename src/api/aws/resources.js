import { call } from "../misc";

export const getEC2 = (token, accountId) => {
  let route = `/ec2?account=${accountId}`;
  return call(route, 'GET', null, token);
};

export const getRDS = (token, accountId) => {
  let route = `/rds?account=${accountId}`;
  return call(route, 'GET', null, token);
};
