import { combineReducers } from 'redux';
import aws from './aws';
import gcp from './gcp';
import auth from './auth';
import user from './user';
import dashboard from './dashboard';
import events from './events';

export default combineReducers({
  aws,
  gcp,
  auth,
  user,
  dashboard,
  events
});
