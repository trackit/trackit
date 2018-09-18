import { combineReducers } from 'redux';
import dates from './datesReducer';
import costs from './costsReducer';
import events from './eventsReducer';

export default combineReducers({
    dates,
    costs,
    events
});
