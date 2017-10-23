import { put, call, all, select } from 'redux-saga/effects';
import API from '../../api';
import Constants from '../../constants';

const getToken = state => state.auth.token;

export function* getAccessSaga() {
  try {
    const token = yield select(getToken);
    const res = yield call(API.AWS.Access.getAccess, token);
    yield all([
      put({ type: Constants.AWS_GET_ACCESS_SUCCESS, access: res.data.access }),
    ]);
  } catch (error) {
    yield put({ type: Constants.AWS_GET_ACCESS_ERROR, error });
  }
}
