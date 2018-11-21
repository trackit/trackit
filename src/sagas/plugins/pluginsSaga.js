import { put, call } from 'redux-saga/effects';
import API from '../../api';
import Constants from '../../constants';
import {getAWSAccounts, getToken} from "../misc";

export function* getPluginsDataSaga() {
  try {
    const token = yield getToken();
    const accounts = yield getAWSAccounts();
    const res = yield call(API.Plugins.getData, token, accounts);
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data") && !res.data.hasOwnProperty("error"))
      yield put({ type: Constants.GET_PLUGINS_DATA_SUCCESS, data: res.data });
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.GET_PLUGINS_DATA_ERROR, error });
  }
}