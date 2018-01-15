import { combineReducers } from 'redux';
import token from './tokenReducer';
import mail from './mailReducer';
import registration from './registrationReducer';
import loginStatus from './loginReducer';

export default combineReducers({
  token,
  mail,
  registration,
  loginStatus
});
