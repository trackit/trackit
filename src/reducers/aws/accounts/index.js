import { combineReducers } from 'redux';
import all from './allReducer';
import external from './externalReducer';
import bills from './billsReducer';

export default combineReducers({
  all,
  external,
  bills
});
