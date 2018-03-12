import { combineReducers } from 'redux';
import items from './itemsReducer';
import values from './valuesReducer';
import dates from './datesReducer';
import intervals from './intervalsReducer';
import filters from './filtersReducer';

export default combineReducers({
  items,
  values,
  dates,
  intervals,
  filters
});
