import { combineReducers } from 'redux';
import keys from './keysReducer';
import selected from './selectedReducer';
import costs from './costsReducer';

export default combineReducers({
    keys,
    selected,
    costs
});
