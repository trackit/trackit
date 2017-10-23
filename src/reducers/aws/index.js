import { combineReducers } from 'redux';
import pricing from './pricingReducer';
import accounts from './accountsReducer';

export default combineReducers({
  pricing,
  accounts
});
