import { put, call } from 'redux-saga/effects';
import Moment from 'moment';
import API from '../../api';
import Constants from '../../constants';
import {getAWSAccounts, getToken} from "../misc";

export function* getEC2ReportSaga({date}) {
  try {
    const token = yield getToken();
    const accounts = yield getAWSAccounts();
    let res;
    if (date.isSameOrAfter(Moment().startOf('months')))
      res = yield call(API.AWS.Resources.getEC2, token, accounts);
    else
      res = yield call(API.AWS.Resources.getEC2History, token, date, accounts);
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
    let res;
    if (date.isSameOrAfter(Moment().startOf('months')))
      res = yield call(API.AWS.Resources.getRDS, token, accounts);
    else {
//      res = yield call(API.AWS.Resources.getRDSHistory, token, date, accounts);
      yield put({ type: Constants.AWS_RESOURCES_GET_RDS_SUCCESS, report: [] });
      return;
    }
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
