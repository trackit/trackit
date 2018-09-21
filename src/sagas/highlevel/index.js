import { takeLatest } from 'redux-saga/effects';
import * as highlevelSaga from './highlevelSaga';
import Constants from '../../constants';

export function* watchHighLevelCostsData() {
  yield takeLatest(Constants.HIGHLEVEL_COSTS_REQUEST, highlevelSaga.getCostsSaga);
}

export function* watchHighLevelEventsData() {
  yield takeLatest(Constants.HIGHLEVEL_EVENTS_REQUEST, highlevelSaga.getEventsSaga);
}

export function* watchHighLevelTagsKeys() {
  yield takeLatest(Constants.HIGHLEVEL_TAGS_KEYS_REQUEST, highlevelSaga.getTagsKeysSaga);
}

export function* watchHighLevelTagsValues() {
  yield takeLatest(Constants.HIGHLEVEL_TAGS_COST_REQUEST, highlevelSaga.getTagsValuesSaga);
}

