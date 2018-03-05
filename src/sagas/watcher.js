import * as AWS from './aws';
import * as GCP from './gcp';
import * as Auth from './auth';
import { takeEvery, takeLatest } from 'redux-saga/effects';
import {getDataSaga, saveDashboardSaga, loadDashboardSaga, initDashboardSaga} from "./dashboardSaga";
import Constants from "../constants";

const Dashboard = {
  watchGetDashboardValues: function*() {
    yield takeEvery(Constants.DASHBOARD_GET_VALUES, getDataSaga);
  },
  watchSaveDashboard: function*() {
    yield takeEvery(Constants.DASHBOARD_UPDATE_ITEMS, saveDashboardSaga);
    yield takeEvery(Constants.DASHBOARD_ADD_ITEM, saveDashboardSaga);
    yield takeEvery(Constants.DASHBOARD_REMOVE_ITEM, saveDashboardSaga);
    yield takeEvery(Constants.DASHBOARD_SET_ITEM_DATES, saveDashboardSaga);
    yield takeEvery(Constants.DASHBOARD_RESET_ITEMS_DATES, saveDashboardSaga);
    yield takeEvery(Constants.DASHBOARD_SET_ITEM_INTERVAL, saveDashboardSaga);
    yield takeEvery(Constants.DASHBOARD_RESET_ITEMS_INTERVAL, saveDashboardSaga);
    yield takeEvery(Constants.DASHBOARD_SET_ITEM_FILTER, saveDashboardSaga);
    yield takeEvery(Constants.DASHBOARD_RESET_ITEMS_FILTER, saveDashboardSaga);
  },
  watchLoadDashboard: function*() {
    yield takeLatest(Constants.DASHBOARD_LOAD_ITEMS, loadDashboardSaga);
  },
  watchInitDashboard: function*() {
    yield takeLatest(Constants.DASHBOARD_INIT_ITEMS, initDashboardSaga);
  }
};

export default {
  ...AWS,
  ...GCP,
  ...Auth,
  ...Dashboard
};
