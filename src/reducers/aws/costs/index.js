import { combineReducers } from 'redux';
import values from './valuesReducer';
import dates from './datesReducer';
import interval from './intervalReducer';

export default combineReducers({
  values,
  dates,
  interval
});
