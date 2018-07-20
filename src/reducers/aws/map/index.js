import { combineReducers } from 'redux';
import values from './valuesReducer';
import dates from './datesReducer';

export default combineReducers({
  values,
  dates,
});
