import { takeLatest } from 'redux-saga/effects';
import * as highlevelSaga from './highlevelSaga';
import Constants from '../../constants';

export function* watchHighLevelCostsData() {
  yield takeLatest(Constants.HIGHLEVEL_COSTS_REQUEST, highlevelSaga.getCostsSaga);
}

export function* watchHighLevelEventsData() {
  yield takeLatest(Constants.HIGHLEVEL_EVENTS_REQUEST, highlevelSaga.getEventsSaga);
}

