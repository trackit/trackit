import { takeLatest } from 'redux-saga/effects';
import * as PluginsSaga from './pluginsSaga';
import Constants from '../../constants';

export function* watchGetPluginsData() {
  yield takeLatest(Constants.GET_PLUGINS_DATA, PluginsSaga.getPluginsDataSaga);
}
