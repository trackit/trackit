import { call } from "../misc";

export const getEC2 = (token, accountId) => {
  let route = `/ec2?aa=${accountId}`;
  return call(route, 'GET', null, token);
};

export const getRDS = (token, accountId) => {
  let route = `/rds?aa=${accountId}`;
  return call(route, 'GET', null, token);
};
