import { takeLatest } from 'redux-saga/effects';
import * as EventsSaga from './eventsSaga';
import * as FiltersSaga from './filtersSaga';
import Constants from '../../constants';

export function* watchGetEventsData() {
  yield takeLatest(Constants.GET_EVENTS_DATA, EventsSaga.getEventsDataSaga);
}

export function* watchGetEventsFilters() {
  yield takeLatest(Constants.EVENTS_GET_FILTERS, FiltersSaga.getEventsFiltersSaga);
}

export function* watchSetEventsFilters() {
  yield takeLatest(Constants.EVENTS_SET_FILTERS, FiltersSaga.setEventsFiltersSaga);
}
export function* watchSnoozeEvent() {
  yield takeLatest(Constants.SNOOZE_EVENT, EventsSaga.snoozeEventSaga);
}

export function* watchUnsnoozeEvent() {
  yield takeLatest(Constants.UNSNOOZE_EVENT, EventsSaga.unsnoozeEventSaga);
}

