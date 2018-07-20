import { put, all } from 'redux-saga/effects';
import Constants from '../../constants';

export default function* logoutSaga() {
  yield all([
    put({ type: Constants.LOGOUT_REQUEST_SUCCESS }),
    put({ type: Constants.CLEAN_USER_TOKEN }),
    put({ type: Constants.CLEAN_USER_MAIL }),
    put({ type: Constants.CLEAN_USER_SELECTED_ACCOUNTS }),
  ]);
}
