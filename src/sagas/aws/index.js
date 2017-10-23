import { takeLatest } from 'redux-saga/effects';
import { getAwsPricingSaga } from './pricingSaga';
import { getAccessSaga } from './accessSaga';

import Constants from '../../constants';

export function* watchGetAwsPricing() {
  yield takeLatest(Constants.AWS_GET_PRICING, getAwsPricingSaga);
}

export function* watchGetAccess() {
  yield takeLatest(Constants.AWS_GET_ACCESS, getAccessSaga);
}