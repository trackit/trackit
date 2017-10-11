import { combineReducers } from 'redux';
import types from './typesReducer';
import aws from './awsReducer';
import gcp from './gcpReducer';


// Combines all reducers to a single reducer function
const rootReducer = combineReducers({
  types,
  aws,
  gcp,
});

export default rootReducer;
