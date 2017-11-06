import { combineReducers } from 'redux';
import accounts from './accounts';
import s3 from './s3Reducer';

export default combineReducers({
  accounts,
  s3,
});
