import { combineReducers } from 'redux';
import aws from './aws';
import gcp from './gcp';
import auth from './auth';
import dashboard from './dashboard';

export default combineReducers({
  aws,
  gcp,
  auth,
  dashboard
});
