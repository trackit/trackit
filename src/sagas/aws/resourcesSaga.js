import { put, call } from 'redux-saga/effects';
import API from '../../api';
import Constants from '../../constants';
import {getAWSAccounts, getToken} from "../misc";

export function* getEC2ReportSaga({date}) {
  try {
    const token = yield getToken();
    const accounts = yield getAWSAccounts();
    const res = yield call(API.AWS.Resources.getEC2, token, date, accounts);
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

export function* getRDSReportSaga({date}) {
  try {
    const token = yield getToken();
    const accounts = yield getAWSAccounts();
    const res = yield call(API.AWS.Resources.getRDS, token, date, accounts);
    yield put({ type: Constants.AWS_RESOURCES_GET_RDS_SUCCESS, report: [] });
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

export function* getESReportSaga({date}) {
  try {
    const token = yield getToken();
    const accounts = yield getAWSAccounts();
    const res = yield call(API.AWS.Resources.getES, token, date, accounts);
    if (res.success && res.hasOwnProperty("data") && !res.data.hasOwnProperty("error"))
      yield put({ type: Constants.AWS_RESOURCES_GET_ES_SUCCESS, report: res.data });
    else if (res.success && res.data.hasOwnProperty("error"))
      throw Error(res.data.error);
    else
      throw Error("Unable to retrieve report");
  } catch (error) {
    yield put({ type: Constants.AWS_RESOURCES_GET_ES_ERROR, error });
  }
}
