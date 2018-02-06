import { put, call } from 'redux-saga/effects';
import API from '../../api';
import Constants from '../../constants';
import {getAWSAccounts, getToken} from "../misc";

export function* getS3DataSaga({ begin, end }) {
  try {
    const token = yield getToken();
    const accounts = yield getAWSAccounts();
    const res = yield call(API.AWS.S3.getData, token, begin, end, accounts);
    if (res.success && res.hasOwnProperty("data") && !res.data.hasOwnProperty("error"))
      yield put({ type: Constants.AWS_GET_S3_DATA_SUCCESS, s3Data: res.data });
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_GET_S3_DATA_ERROR, error });
  }
}
