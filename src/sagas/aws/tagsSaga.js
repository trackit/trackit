import {put, call, all} from 'redux-saga/effects';
import {getToken, getAWSAccounts, initialTagsCharts, getTagsCharts} from '../misc';
import API from '../../api';
import Constants from '../../constants';
import {getTagsCharts as getTagsChartsLS, setTagsCharts} from "../../common/localStorage";

export function* getTagsKeysSaga({ id, begin, end }) {
  try {
    const token = yield getToken();
    const accounts = yield getAWSAccounts();
    const res = yield call(API.AWS.Costs.getTagsKeys, token, begin, end, accounts);
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data")) {
      if (res.data.hasOwnProperty("error"))
        throw Error(res.data.error);
      else
        yield put({type: Constants.AWS_TAGS_GET_KEYS_SUCCESS, tags: res.data, id});
    }
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({type: Constants.AWS_TAGS_GET_KEYS_ERROR, error});
  }
}

export function* getTagsValuesSaga({ id, begin, end, key }) {
  try {
    const token = yield getToken();
    const accounts = yield getAWSAccounts();
    const res = yield call(API.AWS.Costs.getTagsValues, token, begin, end, key, ["product"], accounts);
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data")) {
      if (res.data.hasOwnProperty("error"))
        throw Error(res.data.error);
      else if (res.data.hasOwnProperty(key) && Array.isArray(res.data[key]))
        yield put({type: Constants.AWS_TAGS_GET_VALUES_SUCCESS, values: res.data[key], id});
      else
        throw Error("Error with response");
    }
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({type: Constants.AWS_TAGS_GET_VALUES_ERROR, error});
  }
}

export function* initTagsChartsSaga() {
  try {
    const data = yield call(initialTagsCharts);
    if (data.hasOwnProperty("charts") && data.hasOwnProperty("dates") && data.hasOwnProperty("interval")) {
      yield all([
        put({type: Constants.AWS_TAGS_INSERT_CHARTS, charts: data.charts}),
        put({type: Constants.AWS_TAGS_INSERT_DATES, dates: data.dates}),
        put({type: Constants.AWS_TAGS_INSERT_INTERVAL, interval: data.interval})
      ]);
      setTagsCharts(data);
    }
    else
      throw Error("Invalid data for cost breakdown charts");
  } catch (error) {
    yield put({ type: Constants.AWS_TAGS_INIT_CHARTS_ERROR, error });
  }
}

export function* loadTagsChartsSaga() {
  try {
    const data = yield call(getTagsChartsLS);
    if (!data || (data.hasOwnProperty("charts") && Array.isArray(data.charts)))
      throw Error("No tags chart available");
    else if (data.hasOwnProperty("charts") && data.hasOwnProperty("dates") && data.hasOwnProperty("interval"))
      yield all([
        put({type: Constants.AWS_TAGS_INSERT_CHARTS, charts: data.charts}),
        put({type: Constants.AWS_TAGS_INSERT_DATES, dates: data.dates}),
        put({type: Constants.AWS_TAGS_INSERT_INTERVAL, interval: data.interval})
      ]);
    else
      throw Error("Invalid data for tags charts");
  } catch (error) {
    yield put({ type: Constants.AWS_TAGS_LOAD_CHARTS_ERROR, error });
  }
}

export function* saveTagsChartsSaga() {
  const data = yield getTagsCharts();
  setTagsCharts(data);
}
