import { takeLatest } from 'redux-saga/effects';
import { getAwsPricingSaga } from './pricingSaga';

import Constants from '../../constants';

export function* watchGetAwsPricing() {
  yield takeLatest(Constants.AWS_GET_PRICING, getAwsPricingSaga);
}
