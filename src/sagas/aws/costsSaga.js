import { all, put, call } from 'redux-saga/effects';
import { getToken, getAWSAccounts, getCostBreakdownCharts } from '../misc';
import { setCostBreakdownCharts, getCostBreakdownCharts as getCostBreakdownChartsLS } from '../../common/localStorage';
import UUID from 'uuid/v4';
import moment from 'moment/moment'
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
    else if (data.hasOwnProperty("charts") && data.hasOwnProperty("dates") && data.hasOwnProperty("interval") && data.hasOwnProperty("filter")) {
      yield all([
        put({type: Constants.AWS_INSERT_CHARTS, charts: data.charts}),
        put({type: Constants.AWS_INSERT_COSTS_DATES, dates: data.dates}),
        put({type: Constants.AWS_INSERT_COSTS_INTERVAL, interval: data.interval}),
        put({type: Constants.AWS_INSERT_COSTS_FILTER, filter: data.filter})
      ]);
    } else
      throw Error("Invalid data for cost breakdown charts");
  } catch (error) {
    yield put({ type: Constants.AWS_LOAD_CHARTS_ERROR, error });
  }
}

export function* initChartsSaga() {
  try {
    const initialCharts = [UUID(), UUID()];
    let initialDates = {};
    initialCharts.forEach((id) => {
      initialDates[id] = {
        startDate: moment().subtract(1, 'month').startOf('month'),
        endDate: moment().subtract(1, 'month').endOf('month')
      };
    });
    let initialIntervals = {};
    initialIntervals[initialCharts[0]] = "day";
    initialIntervals[initialCharts[1]] = "week";
    let initialFilters = {};
    initialFilters[initialCharts[0]] = "product";
    initialFilters[initialCharts[1]] = "region";
    yield all([
      put({type: Constants.AWS_INSERT_CHARTS, charts: initialCharts}),
      put({type: Constants.AWS_INSERT_COSTS_DATES, dates: initialDates}),
      put({type: Constants.AWS_INSERT_COSTS_INTERVAL, interval: initialIntervals}),
      put({type: Constants.AWS_INSERT_COSTS_FILTER, filter: initialFilters})
    ]);
  } catch (error) {
    yield put({ type: Constants.AWS_INIT_CHARTS_ERROR, error });
  }
}