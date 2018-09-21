import { combineReducers } from 'redux';
import dates from './datesReducer';
import costs from './costsReducer';
import events from './eventsReducer';
import tags from './tags';

export default combineReducers({
    dates,
    costs,
    events,
    tags
});
