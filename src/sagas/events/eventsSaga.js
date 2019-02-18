import { put, call } from 'redux-saga/effects';
import API from '../../api';
import Constants from '../../constants';
import {getAWSAccounts, getToken, getEventDates} from "../misc";

export function* getEventsDataSaga({ begin, end }) {
  try {
    const token = yield getToken();
    const accounts = yield getAWSAccounts();
    const res = yield call(API.Events.getData, token, begin, end, accounts);
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data") && !res.data.hasOwnProperty("error"))
      yield put({ type: Constants.GET_EVENTS_DATA_SUCCESS, data: res.data });
    else if (res.success && res.hasOwnProperty("data") && res.data.hasOwnProperty("error"))
      throw Error(res.data.error);
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.GET_EVENTS_DATA_ERROR, error });
  }
}

export function* snoozeEventSaga({id}) {
  const token = yield getToken();
  const dates = yield getEventDates();
  const res = yield call(API.Events.snoozeEvent, token, id);

  if (res.success) {
    yield put({
      type: Constants.GET_EVENTS_DATA,
      begin: dates.startDate,
      end: dates.endDate
    });
  }
}

export function* unsnoozeEventSaga({id}) {
  const token = yield getToken();
  const dates = yield getEventDates();
  const res = yield call(API.Events.unsnoozeEvent, token, id);

  if (res.success) {
    yield put({
      type: Constants.GET_EVENTS_DATA,
      begin: dates.startDate,
      end: dates.endDate
    });
  }
}
