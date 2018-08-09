import { combineReducers } from 'redux';
import charts from './chartsReducer';
import keys from './keysReducer';
import dates from './datesReducer';
import values from './valuesReducer';
import interval from './intervalReducer';

export default combineReducers({
  charts,
  keys,
  dates,
  values,
  interval
});
