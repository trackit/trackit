import { combineReducers } from 'redux';
import aws from './aws';
import gcp from './gcp';
import auth from './auth';


export default combineReducers({
  aws,
  gcp,
  auth,
});
