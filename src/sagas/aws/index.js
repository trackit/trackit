import { takeLatest } from 'redux-saga/effects';
import { getAwsPricingSaga } from './pricingSaga';
import { getS3DataSaga } from './s3Saga';

import Constants from '../../constants';

export function* watchGetAwsPricing() {
  yield takeLatest(Constants.AWS_GET_PRICING, getAwsPricingSaga);
}

export function* watchGetAwsS3Data() {
  yield takeLatest(Constants.AWS_GET_S3_DATA, getS3DataSaga);
  yield takeLatest(Constants.AWS_SET_S3_VIEW_DATES, getS3DataSaga);
}
