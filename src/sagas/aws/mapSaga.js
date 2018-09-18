import { put, call } from 'redux-saga/effects';
import { getToken, getAWSAccounts } from '../misc';
import API from '../../api';
import Constants from '../../constants';

export function* getMapCostsSaga({ begin, end, filter }) {
  try {
    const token = yield getToken();
    const accounts = yield getAWSAccounts();
    const res = yield call(API.AWS.Costs.getCosts, token, begin, end, [filter, "product"], accounts);
    if (res.success === null) {
      yield put({type: Constants.LOGOUT_REQUEST});
      return;
    }
    if (res.success && res.hasOwnProperty("data")) {
      if (res.data.hasOwnProperty("error"))
        throw Error(res.data.error);
      else
        yield put({type: Constants.AWS_MAP_GET_COSTS_SUCCESS, costs: res.data});
    }
    else
      throw Error("Error with request");
  } catch (error) {
    yield put({type: Constants.AWS_MAP_GET_COSTS_ERROR, error});
  }
}
