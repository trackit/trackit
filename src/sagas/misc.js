import { select } from 'redux-saga/effects';
import Validations from '../common/forms/AWSAccountForm';
import moment from "moment/moment";
import UUID from "uuid/v4";

const getAccountIDFromRole = Validations.getAccountIDFromRole;

const getTokenFromState = (state) => (state.auth.token);

export const getToken = () => {
  return select(getTokenFromState);
};

const getAWSAccountsFromState = (state) => (state.aws.accounts.selection.map((account) => (getAccountIDFromRole(account.roleArn))));

export const getAWSAccounts = () => {
  return select(getAWSAccountsFromState);
};

const getCostBreakdownChartsFromState = (state) => {
  let data = Object.assign({}, state.aws.costs);
  delete data.values;
  return data;
};

export const getCostBreakdownCharts = () => {
  return select(getCostBreakdownChartsFromState)
};

export const initialCostBreakdownCharts = () => {
  const charts = [UUID(), UUID()];
  let dates = {};
  charts.forEach((id) => {
    dates[id] = {
      startDate: moment().subtract(1, 'month').startOf('month'),
      endDate: moment().subtract(1, 'month').endOf('month')
    };
  });
  let interval = {};
  interval[charts[0]] = "day";
  interval[charts[1]] = "week";
  let filter = {};
  filter[charts[0]] = "product";
  filter[charts[1]] = "region";
  return { charts, dates, interval, filter };
};
