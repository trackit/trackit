import { combineReducers } from 'redux';
import all from './allReducer';
import external from './externalReducer';

export default combineReducers({
  all,
  external,
});
