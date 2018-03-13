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

const getDashboardFromState = (state) => {
  let data = Object.assign({}, state.dashboard);
  delete data.values;
  return data;
};

export const getDashboard = () => {
  return select(getDashboardFromState);
};

export const initialCostBreakdownCharts = () => {
  const id1 = UUID();
  const id2 = UUID();
  const id3 = UUID();
  let charts = {};
  charts[id1] = "bar";
  charts[id2] = "pie";
  charts[id3] = "bar";
  let dates = {};
  Object.keys(charts).forEach((id) => {
    dates[id] = {
      startDate: moment().subtract(1, 'month').startOf('month'),
      endDate: moment().subtract(1, 'month').endOf('month')
    };
  });
  let interval = {};
  interval[id1] = "day";
  interval[id2] = "month";
  interval[id3] = "day";
  let filter = {};
  filter[id1] = "product";
  filter[id2] = "product";
  filter[id3] = "region";
  return { charts, dates, interval, filter };
};

export const initialDashboard = () => {
  const id1 = UUID();
  const id2 = UUID();
  const id3 = UUID();
  const id4 = UUID();
  const id5 = UUID();
  const id6 = UUID();
  const id7 = UUID();
  let items = {};
  items[id1] = {maxSize:[6, null], position:[0, 1], size:[6, 2], static: false, type: "cb_infos"};
  items[id2] = {maxSize:[6, null], position:[0, 3], size:[3, 4], static: false, type: "cb_bar"};
  items[id3] = {maxSize:[6, null], position:[3, 3], size:[3, 4], static: false, type: "cb_pie"};
  items[id4] = {maxSize:[6, null], position:[0, 7], size:[6, 2], static: false, type: "s3_infos"};
  items[id5] = {maxSize:[6, null], position:[0, 9], size:[2, 4], static: false, type: "s3_chart"};
  items[id6] = {maxSize:[6, null], position:[2, 9], size:[2, 4], static: false, type: "s3_chart"};
  items[id7] = {maxSize:[6, null], position:[4, 9], size:[2, 4], static: false, type: "s3_chart"};
  let dates = {
    startDate: moment().subtract(1, 'month').startOf('month'),
    endDate: moment().subtract(1, 'month').endOf('month')
  };
  let intervals = {};
  intervals[id1] = "month";
  intervals[id2] = "day";
  intervals[id3] = "month";
  intervals[id4] = "day";
  intervals[id5] = "day";
  intervals[id6] = "day";
  intervals[id7] = "day";
  let filters = {};
  filters[id1] = null;
  filters[id2] = "product";
  filters[id3] = "product";
  filters[id4] = null;
  filters[id5] = "storage";
  filters[id6] = "requests";
  filters[id7] = "bandwidth";
  return { items, dates, intervals, filters };
};
