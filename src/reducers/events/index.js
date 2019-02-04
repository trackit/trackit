import { combineReducers } from 'redux';
import dates from './datesReducer';
import values from './valuesReducer';
import getFilters from './getFiltersReducer';
import setFilters from './setFiltersReducer';

export default combineReducers({
    dates,
    values,
    getFilters,
    setFilters
});
