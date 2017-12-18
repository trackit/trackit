import { combineReducers } from 'redux';
import accounts from './accounts';
import s3 from './s3Reducer';
import costs from './costsReducer';

export default combineReducers({
  accounts,
  s3,
  costs
});
