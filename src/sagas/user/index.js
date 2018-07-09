import { takeLatest } from 'redux-saga/effects';
import * as ViewersSaga from './viewersSaga';
import Constants from '../../constants';

export function* watchNewViewer() {
  yield takeLatest(Constants.USER_NEW_VIEWER, ViewersSaga.newViewerSaga);
}

export function* watchGetViewers() {
  yield takeLatest(Constants.USER_GET_VIEWERS, ViewersSaga.getViewersSaga);
}
