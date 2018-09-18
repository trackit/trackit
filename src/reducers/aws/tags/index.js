import { combineReducers } from 'redux';
import charts from './chartsReducer';
import keys from './keysReducer';
import dates from './datesReducer';
import values from './valuesReducer';
import filters from './filtersReducer';

export default combineReducers({
  charts,
  keys,
  dates,
  values,
  filters
});
