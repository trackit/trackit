import { put, call, all } from 'redux-saga/effects';
import { setToken } from '../../common/localStorage';
import API from '../../api';
import Constants from '../../constants';

export default function* loginSaga({ username, password }) {
  try {
    const res = yield call(API.Auth.login, username, password);
    if (res.success && res.hasOwnProperty("data") && res.data.hasOwnProperty("token")) {
      setToken(res.data.token);
      yield all([
        put({type: Constants.LOGIN_REQUEST_SUCCESS}),
        put({type: Constants.GET_USER_TOKEN}),
      ]);
    }
    else if (res.success && res.hasOwnProperty("data") && res.data.hasOwnProperty("error"))
      throw Error(res.data.error);
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.LOGIN_REQUEST_ERROR, error: error.message });
  }
}
