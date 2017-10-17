import { put, call, all } from 'redux-saga/effects';
import { getToken, setToken } from '../../common/localStorage';
import API from '../../api';
import Constants from '../../constants';


export function* getUserTokenSaga() {
  try {
    const token = yield call(getToken);
    yield all([
      put({ type: Constants.GET_USER_TOKEN_SUCCESS, token }),
    ]);
  } catch (error) {
    yield put({ type: Constants.GET_USER_TOKEN_ERROR, error });
  }
}

export function* loginSaga({ username, password }) {
  try {
    const res = yield call(API.Auth.login, username, password);
    if (res.success && res.token) {
      setToken(res.token);
    }
    yield all([
      put({ type: Constants.LOGIN_REQUEST_SUCCESS }),
      put({ type: Constants.GET_USER_TOKEN }),
    ]);
  } catch (error) {
    yield put({ type: Constants.LOGIN_REQUEST_ERROR, error });
  }
}
