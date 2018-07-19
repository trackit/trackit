import { combineReducers } from 'redux';
import all from './allReducer';
import creation from './creationReducer';

export default combineReducers({
  all,
  creation,
});
