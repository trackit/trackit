import { combineReducers } from 'redux';
import dates from './datesReducer';
import values from './valuesReducer';

export default combineReducers({
    dates,
    values
});
