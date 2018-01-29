import { combineReducers } from 'redux';
import charts from './chartsReducer';
import values from './valuesReducer';
import dates from './datesReducer';
import interval from './intervalReducer';
import filter from './filterReducer';

export default combineReducers({
  charts,
  values,
  dates,
  interval,
  filter
});
