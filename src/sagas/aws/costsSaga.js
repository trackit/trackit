import { all, put, call } from 'redux-saga/effects';
import { getToken, getAWSAccounts, getCostBreakdownCharts, initialCostBreakdownCharts } from '../misc';
import { setCostBreakdownCharts, getCostBreakdownCharts as getCostBreakdownChartsLS } from '../../common/localStorage';
import API from '../../api';
import Constants from '../../constants';

export function* getCostsSaga({ id, begin, end, filters }) {
  try {
    const token = yield getToken();
    const accounts = yield getAWSAccounts();
    const res = yield call(API.AWS.Costs.getCosts, token, begin, end, filters, accounts);
    if (res.success && res.hasOwnProperty("data"))
      yield put({type: Constants.AWS_GET_COSTS_SUCCESS, id, costs: res.data});
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({type: Constants.AWS_GET_COSTS_ERROR, error});
  }
}

export function* saveChartsSaga() {
  const data = yield getCostBreakdownCharts();
  setCostBreakdownCharts(data);
}

export function* loadChartsSaga() {
  try {
    const data = yield call(getCostBreakdownChartsLS);
    if (!data)
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
    if (data.hasOwnProperty("charts") && data.hasOwnProperty("dates") && data.hasOwnProperty("interval") && data.hasOwnProperty("filter"))
      yield all([
        put({type: Constants.AWS_INSERT_CHARTS, charts: data.charts}),
        put({type: Constants.AWS_INSERT_COSTS_DATES, dates: data.dates}),
        put({type: Constants.AWS_INSERT_COSTS_INTERVAL, interval: data.interval}),
        put({type: Constants.AWS_INSERT_COSTS_FILTER, filter: data.filter})
      ]);
    else
      throw Error("Invalid data for cost breakdown charts");
  } catch (error) {
    yield put({ type: Constants.AWS_INIT_CHARTS_ERROR, error });
  }
}