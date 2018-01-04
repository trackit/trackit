import { combineReducers } from 'redux';
import all from './allReducer';
import selection from './selectionReducer';
import external from './externalReducer';
import bills from './billsReducer';

export default combineReducers({
  all,
  selection,
  external,
  bills
});
