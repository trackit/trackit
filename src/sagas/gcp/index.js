import { takeLatest } from 'redux-saga/effects';
import { getGcpPricingSaga } from './pricingSaga';

import Constants from '../../constants';

export function* watchGetGcpPricing() {
  yield takeLatest(Constants.GCP_GET_PRICING, getGcpPricingSaga);
}
