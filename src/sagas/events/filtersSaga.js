import { put, call } from 'redux-saga/effects';
import API from '../../api';
import Constants from '../../constants';
import {getToken} from "../misc";

export function* getEventsFiltersSaga() {
  try {
    const token = yield getToken();
    const res = yield call(API.Events.getFilters, token);
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data") && res.data.hasOwnProperty("filters"))
      yield put({ type: Constants.EVENTS_GET_FILTERS_SUCCESS, data: res.data.filters.map((filter, id) => ({...filter, id})) });
    else if (res.success && res.hasOwnProperty("data") && res.data.hasOwnProperty("error"))
      throw Error(res.data.error);
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.EVENTS_GET_FILTERS_ERROR, error });
  }
}

export function* setEventsFiltersSaga({filters}) {
  try {
    const token = yield getToken();
    const res = yield call(API.Events.setFilters, token, filters);
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data") && res.data.hasOwnProperty("filters"))
      yield put({ type: Constants.EVENTS_SET_FILTERS_SUCCESS, data: res.data.filters });
    else if (res.success && res.hasOwnProperty("data") && res.data.hasOwnProperty("error"))
      throw Error(res.data.error);
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.EVENTS_SET_FILTERS_ERROR, error });
  }
}
