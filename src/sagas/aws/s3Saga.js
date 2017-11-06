import { put, call, all } from 'redux-saga/effects';
import API from '../../api';
import Constants from '../../constants';

export function* getS3DataSaga() {
  try {
    const res = yield call(API.AWS.S3.getS3Data);
    if (res.success && res.hasOwnProperty("data"))
      yield all([
        put({ type: Constants.AWS_GET_S3_DATA_SUCCESS, s3Data: res.data }),
      ]);
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_GET_S3_DATA_ERROR, error });
  }
}
