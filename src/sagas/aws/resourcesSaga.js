import { put, call } from 'redux-saga/effects';
import API from '../../api';
import Constants from '../../constants';
import { getToken } from "../misc";

export function* getEC2ReportSaga({ accountId }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Resources.getEC2, token, accountId);
    if (res.success && res.hasOwnProperty("data") && !res.data.hasOwnProperty("error"))
      yield put({ type: Constants.AWS_RESOURCES_GET_EC2_SUCCESS, report: res.data });
    else if (res.success && res.data.hasOwnProperty("error"))
      throw Error(res.data.error);
    else
      throw Error("Unable to retrieve report");
  } catch (error) {
    yield put({ type: Constants.AWS_RESOURCES_GET_EC2_ERROR, error });
  }
}

export function* getRDSReportSaga({ accountId }) {
  try {
    const token = yield getToken();
    const res = yield call(API.AWS.Resources.getRDS, token, accountId);
    console.log(res);
    if (res.success && res.hasOwnProperty("data") && !res.data.hasOwnProperty("error"))
      yield put({ type: Constants.AWS_RESOURCES_GET_RDS_SUCCESS, report: res.data });
    else if (res.success && res.data.hasOwnProperty("error"))
      throw Error(res.data.error);
    else
      throw Error("Unable to retrieve report");
  } catch (error) {
    yield put({ type: Constants.AWS_RESOURCES_GET_RDS_ERROR, error });
  }
}
