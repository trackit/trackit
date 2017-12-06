import { combineReducers } from 'redux';
import token from './tokenReducer';
import registration from './registrationReducer';

export default combineReducers({
  token,
  registration
});
