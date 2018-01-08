import { combineReducers } from 'redux';
import all from './allReducer';
import external from './externalReducer';
import bills from './billsReducer';
import creation from './creationReducer';

export default combineReducers({
  all,
  external,
  bills,
  creation
});
