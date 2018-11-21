import { combineReducers } from 'redux';
import aws from './aws';
import gcp from './gcp';
import auth from './auth';
import dashboard from './dashboard';
import events from './events';
import plugins from './plugins';
import highlevel from './highlevel';

export default combineReducers({
  aws,
  gcp,
  auth,
  dashboard,
  events,
  plugins,
  highlevel
});
