import { takeLatest } from 'redux-saga/effects';
import * as EventsSaga from './eventsSaga';
import Constants from '../../constants';

export function* watchGetEventsData() {
  yield takeLatest(Constants.GET_EVENTS_DATA, EventsSaga.getEventsDataSaga);
}

