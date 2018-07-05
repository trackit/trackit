import { put, call } from 'redux-saga/effects';
import API from '../../api';
import Constants from '../../constants';

export default function* renewPasswordSaga({ username, password, token }) {
  try {
    yield put({ type: Constants.RENEW_PASSWORD_LOADING });
    const res = yield call(API.Auth.renewPassword, username, password, token);
    if (res.success && res.hasOwnProperty("data") && !res.data.hasOwnProperty("error"))
      yield put({type: Constants.RENEW_PASSWORD_SUCCESS});
    else if (res.success && res.hasOwnProperty("data") && res.data.hasOwnProperty("error"))
      throw Error(res.data.error);
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.RENEW_PASSWORD_ERROR, error: error.message });
  }
}
