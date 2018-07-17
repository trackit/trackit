import { all, put, call } from 'redux-saga/effects';
import { getToken, getAWSAccounts, getDashboard, initialDashboard } from './misc';
import { setDashboard, getDashboard as getDashboardLS } from '../common/localStorage';
import API from '../api';
import Constants from '../constants';

export function* getDataSaga({ id, itemType, begin, end, filters }) {
  try {
    const token = yield getToken();
    const accounts = yield getAWSAccounts();
    let res;
    switch (itemType) {
      case "costbreakdown":
        res = yield call(API.AWS.Costs.getCosts, token, begin, end, filters, accounts);
        break;
      case "s3":
        res = yield call(API.AWS.S3.getData, token, begin, end, accounts);
        break;
      default:
        res = null;
    }
    if (res && res.success && res.hasOwnProperty("data")) {
      if (res.data.hasOwnProperty("error"))
        throw Error(res.data.error);
      else
        yield put({type: Constants.DASHBOARD_GET_VALUES_SUCCESS, id, data: res.data});
    }
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({type: Constants.DASHBOARD_GET_VALUES_ERROR, id, error});
  }
}

export function* saveDashboardSaga() {
  const data = yield getDashboard();
  setDashboard(data);
}

export function* loadDashboardSaga() {
  try {
    const data = yield call(getDashboardLS);
    if (!data || (data.hasOwnProperty("items") && Array.isArray(data.items)))
      throw Error("No dashboard available");
    else if (data.hasOwnProperty("items") && data.hasOwnProperty("dates") && data.hasOwnProperty("intervals") && data.hasOwnProperty("filters"))
      yield all([
        put({type: Constants.DASHBOARD_INSERT_ITEMS, items: data.items}),
        put({type: Constants.DASHBOARD_INSERT_DATES, dates: data.dates}),
        put({type: Constants.DASHBOARD_INSERT_ITEMS_INTERVAL, intervals: data.intervals}),
        put({type: Constants.DASHBOARD_INSERT_ITEMS_FILTER, filters: data.filters})
      ]);
    else
      throw Error("Invalid data for dashboard");
  } catch (error) {
    yield put({ type: Constants.DASHBOARD_INIT_ITEMS_ERROR, error });
  }
}

export function* initDashboardSaga() {
  try {
    const data = yield call(initialDashboard);
    if (data.hasOwnProperty("items") && data.hasOwnProperty("dates") && data.hasOwnProperty("intervals") && data.hasOwnProperty("filters")) {
      yield all([
        put({type: Constants.DASHBOARD_INSERT_ITEMS, items: data.items}),
        put({type: Constants.DASHBOARD_INSERT_DATES, dates: data.dates}),
        put({type: Constants.DASHBOARD_INSERT_ITEMS_INTERVAL, intervals: data.intervals}),
        put({type: Constants.DASHBOARD_INSERT_ITEMS_FILTER, filters: data.filters})
      ]);
      setDashboard(data);
    }
    else
      throw Error("Invalid data for dashboard");
  } catch (error) {
    yield put({ type: Constants.DASHBOARD_INIT_ITEMS_ERROR, error });
  }
}