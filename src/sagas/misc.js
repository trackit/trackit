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

const getSelectedAccountsFromState = (state) => (state.aws.accounts.selection);

export const getSelectedAccounts = () => {
  return select(getSelectedAccountsFromState);
};

const getCostBreakdownChartsFromState = (state) => {
  let data = Object.assign({}, state.aws.costs);
  delete data.values;
  return data;
};

export const getCostBreakdownCharts = () => {
  return select(getCostBreakdownChartsFromState)
};

const getS3DatesFromState = (state) => (state.aws.s3.dates);

export const getS3Dates = () => {
  return select(getS3DatesFromState);
};

export const initialCostBreakdownCharts = () => {
  const id1 = UUID();
  const id2 = UUID();
  let charts = {};
  charts[id1] = "bar";
  charts[id2] = "bar";
  let dates = {};
  Object.keys(charts).forEach((id) => {
    dates[id] = {
      startDate: moment().subtract(1, 'month').startOf('month'),
      endDate: moment().subtract(1, 'month').endOf('month')
    };
  });
  let interval = {};
  interval[id1] = "day";
  interval[id2] = "week";
  let filter = {};
  filter[id1] = "product";
  filter[id2] = "region";
  return { charts, dates, interval, filter };
};
