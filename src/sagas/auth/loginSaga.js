import { put, call, all } from 'redux-saga/effects';
import { setToken } from '../../common/localStorage';
import API from '../../api';
import Constants from '../../constants';

export default function* loginSaga({ username, password }) {
  try {
    const res = yield call(API.Auth.login, username, password);
    if (res.success && res.data.token)
      setToken(res.token);
    yield all([
      put({ type: Constants.LOGIN_REQUEST_SUCCESS }),
      put({ type: Constants.GET_USER_TOKEN }),
    ]);
  } catch (error) {
    yield put({ type: Constants.LOGIN_REQUEST_ERROR, error });
  }
}
