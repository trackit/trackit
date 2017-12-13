import { combineReducers } from 'redux';
import token from './tokenReducer';
import registration from './registrationReducer';
import loginStatus from './loginReducer';

export default combineReducers({
  token,
  registration,
  loginStatus
});
