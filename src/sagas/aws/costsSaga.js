import { put, call } from 'redux-saga/effects';
import { getToken } from '../misc';
import API from '../../api';
import Constants from '../../constants';

export function* getCostsSaga({ begin, end, filters, accounts=undefined }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Costs.getCosts, token, begin, end, filters, accounts);
    if (res.success && res.hasOwnProperty("data"))
      yield put({ type: Constants.AWS_GET_COSTS_SUCCESS, costs: res.data });
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({ type: Constants.AWS_GET_COSTS_ERROR, error });
  }
}
