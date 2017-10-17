import { takeLatest } from 'redux-saga/effects';
import { loginSaga, getUserTokenSaga } from './authSaga';

import Constants from '../../constants';

export function* watchGetLogin() {
  yield takeLatest(Constants.LOGIN_REQUEST, loginSaga);
}

export function* watchGetToken() {
  yield takeLatest(Constants.GET_USER_TOKEN, getUserTokenSaga);
}
