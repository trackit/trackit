import { put, all } from 'redux-saga/effects';
import { unsetToken } from '../../common/localStorage';
import Constants from '../../constants';

export default function* logoutSaga() {
  try {
    unsetToken();
    yield all([
      put({ type: Constants.LOGOUT_REQUEST_SUCCESS }),
      put({ type: Constants.CLEAN_USER_TOKEN }),
    ]);
  } catch (error) {
    yield put({ type: Constants.LOGOUT_REQUEST_ERROR, error });
  }
}
