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
      yield put({ type: Constants.HIGHLEVEL_COSTS_SUCCESS, months: res.data.month, history: resHistory.data.month });
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
    if (
      res.success && res.hasOwnProperty("data") && !res.data.hasOwnProperty("error")
    )
      yield put({ type: Constants.HIGHLEVEL_EVENTS_SUCCESS, events: res.data });
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.HIGHLEVEL_EVENTS_ERROR, error });
  }
}