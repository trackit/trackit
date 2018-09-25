import * as AWS from './aws';
import * as GCP from './gcp';
import * as Auth from './auth';
import * as User from './user';
import * as Events from './events';
import * as Highlevel from './highlevel';
import { takeEvery, takeLatest, fork, cancel } from 'redux-saga/effects';
import {getDataSaga, saveDashboardSaga, loadDashboardSaga, initDashboardSaga} from "./dashboardSaga";
import Constants from "../constants";

// To manage concurrency when multiple calls are fired for the same id
let tasks = {};

function* accumulateGetValuesSaga(action) {
  const {id, type} = action;
  if (!tasks[type]) tasks[type] = {};

  if (tasks[type][id]) {
    yield cancel(tasks[type][id]);
  }
  tasks[type][id] = yield fork(getDataSaga, action);
}


const Dashboard = {
  watchGetDashboardValues: function*() {
    yield takeEvery(Constants.DASHBOARD_GET_VALUES, accumulateGetValuesSaga);
  },
  watchSaveDashboard: function*() {
    yield takeEvery(Constants.DASHBOARD_UPDATE_ITEMS, saveDashboardSaga);
    yield takeEvery(Constants.DASHBOARD_ADD_ITEM, saveDashboardSaga);
    yield takeEvery(Constants.DASHBOARD_REMOVE_ITEM, saveDashboardSaga);
    yield takeEvery(Constants.DASHBOARD_SET_DATES, saveDashboardSaga);
    yield takeEvery(Constants.DASHBOARD_RESET_DATES, saveDashboardSaga);
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
  ...Dashboard,
  ...User,
  ...Events,
  ...Highlevel
};
