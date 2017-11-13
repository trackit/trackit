import { combineReducers } from 'redux';
import pricing from './pricingReducer';
import accounts from './accounts';
import s3 from './s3Reducer';

export default combineReducers({
  pricing,
  accounts,
  s3,
});
