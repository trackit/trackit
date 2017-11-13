import { put, call, all } from 'redux-saga/effects';
import API from '../../api';
import Constants from '../../constants';

export function* getS3DataSaga() {
  try {
    const s3Data = yield call(API.AWS.S3.getS3Data);
    yield all([
      put({ type: Constants.AWS_GET_S3_DATA_SUCCESS, s3Data }),
    ]);
  } catch (error) {
    yield put({ type: Constants.AWS_GET_S3_DATA_ERROR, error });
  }
}
