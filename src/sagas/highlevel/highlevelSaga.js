import { put, call } from 'redux-saga/effects';
import moment from 'moment';
import API from '../../api';
import Constants from '../../constants';
import {getAWSAccounts, getToken} from "../misc";

export function* getCostsSaga({ begin, end }) {
  try {
    const token = yield getToken();
    const accounts = yield getAWSAccounts();
    const res = yield call(API.AWS.Costs.getCosts, token, moment(begin).subtract(1, 'months'), end, ['month', 'product'], accounts);
    const historyEndDate = moment(end).month() === moment().month() ?  moment(end).subtract(1, 'months').endOf('month') :  moment(end).endOf('month');
    const resHistory = yield call(API.AWS.Costs.getCosts, token, moment(begin).subtract(12, 'months'), historyEndDate, ['month'], accounts);
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data") && !res.data.hasOwnProperty("error")
      && resHistory.success && resHistory.hasOwnProperty("data") && !resHistory.data.hasOwnProperty("error"))
      yield put({ type: Constants.HIGHLEVEL_COSTS_SUCCESS, months: (res.data.month || {}), history: (resHistory.data.month || {}) });
    else if (res.success && res.hasOwnProperty("data") && res.data.hasOwnProperty("error"))
      throw Error(res.data.error);
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.HIGHLEVEL_COSTS_ERROR, error });
  }
}

export function* getEventsSaga({ begin, end }) {
  try {
    const token = yield getToken();
    const accounts = yield getAWSAccounts();
    const res = yield call(API.Events.getData, token, begin, end, accounts);
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data") && !res.data.hasOwnProperty("error"))
      yield put({ type: Constants.HIGHLEVEL_EVENTS_SUCCESS, events: res.data });
    else if (res.success && res.hasOwnProperty("data") && res.data.hasOwnProperty("error"))
      throw Error(res.data.error);
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.HIGHLEVEL_EVENTS_ERROR, error });
  }
}

export function* getTagsKeysSaga({ begin, end }) {
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
      else {
        yield put({type: Constants.HIGHLEVEL_TAGS_KEYS_SUCCESS, keys: res.data});
        if (res.data.length) {
          // Setting first key as selected
          yield put({type: Constants.HIGHLEVEL_TAGS_KEYS_SELECT, key: res.data[0]});
          // Retrieving values for this key
          yield put({type: Constants.HIGHLEVEL_TAGS_COST_REQUEST, begin, end, key: res.data[0]});
        } else {
          // No keys so no possible selection or costs
          yield put({type: Constants.HIGHLEVEL_TAGS_COST_CLEAR});
          yield put({type: Constants.HIGHLEVEL_TAGS_KEYS_CLEAR_SELECTED});
        }
      }
    }
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({type: Constants.HIGHLEVEL_TAGS_KEYS_ERROR, error});
  }
}

export function* getTagsValuesSaga({ begin, end, key }) {
  try {
    const token = yield getToken();
    const accounts = yield getAWSAccounts();
    const res = yield call(API.AWS.Costs.getTagsValues, token, moment(begin).subtract(1, 'months'), end, key, ['month'], accounts);
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data")) {
      if (res.data.hasOwnProperty("error"))
        throw Error(res.data.error);
      else if (res.data.hasOwnProperty(key) && Array.isArray(res.data[key])) {
        yield put({type: Constants.HIGHLEVEL_TAGS_COST_SUCCESS, values: res.data[key]});
      }
      else
        throw Error("Error with response");
    }
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({type: Constants.HIGHLEVEL_TAGS_COST_ERROR, error});
  }
}

export function* getUnusedEC2Saga({date}) {
  try {
    const token = yield getToken();
    const accounts = yield getAWSAccounts();
    const res = yield call(API.AWS.Resources.getUnusedEC2, token, date, accounts);
    if (res.success && res.hasOwnProperty("data") && !res.data.hasOwnProperty("error"))
      yield put({ type: Constants.HIGHLEVEL_UNUSED_EC2_SUCCESS, data: res.data });
    else if (res.success && res.data.hasOwnProperty("error"))
      throw Error(res.data.error);
    else
      throw Error("Unable to retrieve report");
  } catch (error) {
    yield put({ type: Constants.HIGHLEVEL_UNUSED_EC2_ERROR, error });
  }
}


