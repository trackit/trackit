import { put, call } from 'redux-saga/effects';
import API from '../../api';
import Constants from '../../constants';

export default function* registrationSaga({ username, password }) {
  try {
    yield put({ type: Constants.REGISTRATION_REQUEST_LOADING });
    const res = yield call(API.Auth.register, username, password);
    if (res.success && !res.data.error) {
      yield put({type: Constants.REGISTRATION_SUCCESS, payload: { status: true }});
    }
    else
      throw Error(res.data.error);
  } catch (error) {
    yield put({ type: Constants.REGISTRATION_ERROR, payload: { status: false, error: error.toString() }});
  }
}
