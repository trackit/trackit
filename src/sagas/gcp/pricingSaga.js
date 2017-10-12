import { put, call, all, select } from 'redux-saga/effects';
import API from '../../api';
import Constants from '../../constants';

const getTypes = state => state.types;

export function* getGcpPricingSaga() {
  try {
    const stateTypes = yield select(getTypes);

    const params = {
      frequent: `${stateTypes.frequentValue}${stateTypes.frequentUnit}`,
      infrequent: `${stateTypes.infrequentValue}${stateTypes.infrequentUnit}`,
      archive: `${stateTypes.archiveValue}${stateTypes.archiveUnit}`,
    };

    const pricing = yield call(API.getPricing, 'gcp', params.frequent, params.infrequent, params.archive );
    yield all([
      put({ type: Constants.GCP_GET_PRICING_SUCCESS, pricing }),
    ]);
  } catch (error) {
    yield put({ type: Constants.GCP_GET_PRICING_ERROR, error });
  }
}
