import { put, call, all } from 'redux-saga/effects';
import { getToken } from '../misc';
import API from '../../api';
import Constants from '../../constants';

export function* getViewersSaga() {
  try {
    const token = yield getToken();
    const res = yield call(API.User.Viewers.list, token);
    if (res.success && res.hasOwnProperty('data'))
      yield put({ type: Constants.USER_GET_VIEWERS_SUCCESS, viewers: res.data });
    else
      throw Error('Error with request');
  } catch (error) {
    yield put({ type: Constants.USER_GET_VIEWERS_ERROR, error });
  }
}

export function* newViewerSaga({ email }) {
  try {
    const token = yield getToken();
    const res = yield call(API.User.Viewers.create, email, token);
    if (res.success && res.hasOwnProperty('data'))
      yield put({ type: Constants.USER_NEW_VIEWER_SUCCESS, viewer: res.data });
    else
      throw Error('Error with request');
  } catch (error) {
    yield put({ type: Constants.USER_NEW_VIEWER_ERROR, error });
  }
}
