import { combineReducers } from 'redux';
import pricing from './pricingReducer';
import s3 from './s3Reducer';

export default combineReducers({
  pricing,
  s3,
});
