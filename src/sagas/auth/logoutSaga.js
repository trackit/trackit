import { put, all } from 'redux-saga/effects';
import Constants from '../../constants';

export default function* logoutSaga() {
  yield all([
    put({ type: Constants.LOGOUT_REQUEST_SUCCESS }),
    put({ type: Constants.CLEAN_USER_TOKEN }),
    put({ type: Constants.CLEAN_USER_MAIL }),
    put({ type: Constants.CLEAN_USER_SELECTED_ACCOUNTS }),
    /* Clean dashboard, cost breakdown, s3 analytics, tags */
    put({ type: Constants.DASHBOARD_CLEAN_ITEMS }),
    put({ type: Constants.AWS_CLEAR_CHARTS }),
    put({ type: Constants.AWS_CLEAR_S3_DATES }),
    put({ type: Constants.AWS_TAGS_CLEAN_CHARTS }),
  ]);
}
