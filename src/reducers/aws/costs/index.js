import { combineReducers } from 'redux';
import values from './valuesReducer';
import dates from './datesReducer';
import interval from './intervalReducer';
import filter from './filterReducer';

export default combineReducers({
  values,
  dates,
  interval,
  filter
});
