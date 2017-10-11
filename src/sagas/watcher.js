import { takeLatest } from 'redux-saga/effects';
import { getGcpPricingSaga, getAwsPricingSaga } from './pricingSaga';

import * as types from '../constants/actionTypes';

export function* watchGetGcpPricing() {
  yield takeLatest(types.GET_PRICING_GCP, getGcpPricingSaga);
}

export function* watchGetAwsPricing() {
  yield takeLatest(types.GET_PRICING_AWS, getAwsPricingSaga);
}
