import { put, call, all, select } from 'redux-saga/effects';
import { getPricing } from '../api/api';
import * as types from '../constants/actionTypes';

const getTypes = state => state.types;

export function* getGcpPricingSaga() {
  try {
    const stateTypes = yield select(getTypes);

    const params = {
      frequent: `${stateTypes.frequentValue}${stateTypes.frequentUnit}`,
      infrequent: `${stateTypes.infrequentValue}${stateTypes.infrequentUnit}`,
      archive: `${stateTypes.archiveValue}${stateTypes.archiveUnit}`,
    };

    const pricing = yield call(getPricing, 'gcp', params.frequent, params.infrequent, params.archive );
    yield all([
      put({ type: types.GET_PRICING_GCP_SUCCESS, pricing }),
    ]);
  } catch (error) {
    yield put({ type: types.GET_PRICING_GCP_ERROR, error });
  }
}

export function* getAwsPricingSaga() {
  try {
    const stateTypes = yield select(getTypes);

    const params = {
      frequent: `${stateTypes.frequentValue}${stateTypes.frequentUnit}`,
      infrequent: `${stateTypes.infrequentValue}${stateTypes.infrequentUnit}`,
      archive: `${stateTypes.archiveValue}${stateTypes.archiveUnit}`,
    };

    const pricing = yield call(getPricing, 'aws', params.frequent, params.infrequent, params.archive );
    yield all([
      put({ type: types.GET_PRICING_AWS_SUCCESS, pricing }),
    ]);
  } catch (error) {
    yield put({ type: types.GET_PRICING_AWS_ERROR, error });
  }
}
