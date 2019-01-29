import { takeLatest } from 'redux-saga/effects';
import * as EventsSaga from './eventsSaga';
import Constants from '../../constants';

export function* watchGetEventsData() {
  yield takeLatest(Constants.GET_EVENTS_DATA, EventsSaga.getEventsDataSaga);
}

export function* watchSnoozeEvent() {
  yield takeLatest(Constants.SNOOZE_EVENT, EventsSaga.snoozeEventSaga);
}

export function* watchUnsnoozeEvent() {
  yield takeLatest(Constants.UNSNOOZE_EVENT, EventsSaga.unsnoozeEventSaga);
}

