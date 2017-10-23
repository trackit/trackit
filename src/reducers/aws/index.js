import { combineReducers } from 'redux';
import pricing from './pricingReducer';
import access from './accessReducer';

export default combineReducers({
  pricing,
  access
});
