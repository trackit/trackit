import { put, call } from 'redux-saga/effects';
import API from '../../api';
import Constants from '../../constants';

export default function* registrationSaga({ username, password, awsToken }) {
  try {
    yield put({ type: Constants.REGISTRATION_REQUEST_LOADING });
    const res = yield call(API.Auth.register, username, password, awsToken);
    if (res.success && !res.data.error) {
      yield put({type: Constants.REGISTRATION_SUCCESS, payload: { status: true }});
    }
    else if (res.data && res.data.error)
      throw Error(res.data.error);
    else
      throw Error('An error has occured');
  } catch (error) {
    yield put({
      type: Constants.REGISTRATION_ERROR,
      payload: { status: false, error: error.message }
    });
  }
}
