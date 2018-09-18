import { combineReducers } from 'redux';
import values from './valuesReducer';
import dates from './datesReducer';
import filter from './filterReducer';

export default combineReducers({
  values,
  dates,
  filter
});
