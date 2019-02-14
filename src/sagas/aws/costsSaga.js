import { all, put, call } from 'redux-saga/effects';
import { getToken, getAWSAccounts, getCostBreakdownCharts, initialCostBreakdownCharts } from '../misc';
import { setCostBreakdownCharts, getCostBreakdownCharts as getCostBreakdownChartsLS, unsetCostBreakdownCharts } from '../../common/localStorage';
import API from '../../api';
import Constants from '../../constants';
import Validation from '../../common/forms/AWSAccountForm';

const getAccountIDFromRole = Validation.getAccountIDFromRole;

export function* getCostsSaga({ id, begin, end, filters, chartType }) {
  try {
    const token = yield getToken();
    const accounts = yield getAWSAccounts();
    let res;
    if (chartType === "breakdown")
      res = yield call(API.AWS.Costs.getCosts, token, begin, end, filters, accounts);
    else if (chartType === "differentiator")
      res = yield call(API.AWS.Costs.getCostDiff, token, begin, end, filters, accounts);
    else
      res = {success: false};
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data")) {
      if (res.data.hasOwnProperty("error"))
        throw Error(res.data.error);
      else {
        if (res.data.hasOwnProperty("account")) {
          const accountsRaw = yield call(API.AWS.Accounts.getAccounts, token);
          if (accountsRaw.success && accountsRaw.hasOwnProperty("data")) {
            const accounts = {};
            accountsRaw.data.forEach((item) => {
              const accountID = (item.hasOwnProperty("awsIdentity") ? item.awsIdentity : getAccountIDFromRole(item.roleArn));
              accounts[accountID] = {...item, accountID};
              if (item.subAccounts) {
                item.subAccounts.forEach((item) => {
                  accounts[item.awsIdentity] = {...item, accountID: item.awsIdentity};
                });
              }
            });
            const newData = {};
            Object.keys(res.data.account).forEach((accountID) => {
              if (Object.keys(accounts).indexOf(accountID) !== -1)
                newData[accounts[accountID].pretty] = res.data.account[accountID];
              else
                newData[accountID] = res.data.account[accountID];
            });
            res.data.account = newData;
          }
          else
            throw Error("Error while getting accounts");
        }
        yield put({type: Constants.AWS_GET_COSTS_SUCCESS, id, costs: res.data});
      }
    }
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({type: Constants.AWS_GET_COSTS_ERROR, id, error});
  }
}

export function* saveChartsSaga() {
  const data = yield getCostBreakdownCharts();
  setCostBreakdownCharts(data);
}

export function* loadChartsSaga() {
  try {
    const data = yield call(getCostBreakdownChartsLS);
    if (!data || (data.hasOwnProperty("charts") && Array.isArray(data.charts)))
      throw Error("No cost breakdown chart available");
    else if (data.hasOwnProperty("charts") && data.hasOwnProperty("dates") && data.hasOwnProperty("interval") && data.hasOwnProperty("filter"))
      yield all([
        put({type: Constants.AWS_INSERT_CHARTS, charts: data.charts}),
        put({type: Constants.AWS_INSERT_COSTS_DATES, dates: data.dates}),
        put({type: Constants.AWS_INSERT_COSTS_INTERVAL, interval: data.interval}),
        put({type: Constants.AWS_INSERT_COSTS_FILTER, filter: data.filter})
      ]);
    else
      throw Error("Invalid data for cost breakdown charts");
  } catch (error) {
    yield put({ type: Constants.AWS_LOAD_CHARTS_ERROR, error });
  }
}

export function* initChartsSaga() {
  try {
    const data = yield call(initialCostBreakdownCharts);
    if (data.hasOwnProperty("charts") && data.hasOwnProperty("dates") && data.hasOwnProperty("interval") && data.hasOwnProperty("filter")) {
      yield all([
        put({type: Constants.AWS_INSERT_CHARTS, charts: data.charts}),
        put({type: Constants.AWS_INSERT_COSTS_DATES, dates: data.dates}),
        put({type: Constants.AWS_INSERT_COSTS_INTERVAL, interval: data.interval}),
        put({type: Constants.AWS_INSERT_COSTS_FILTER, filter: data.filter})
      ]);
      setCostBreakdownCharts(data);
    }
    else
      throw Error("Invalid data for cost breakdown charts");
  } catch (error) {
    yield put({ type: Constants.AWS_INIT_CHARTS_ERROR, error });
  }
}

export function* clearChartsSaga() {
  yield call(unsetCostBreakdownCharts);
}
