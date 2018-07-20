import { combineReducers } from 'redux';
import token from './tokenReducer';
import mail from './mailReducer';
import registration from './registrationReducer';
import loginStatus from './loginReducer';
import recoverStatus from './recoverReducer';
import renewStatus from './renewReducer';

export default combineReducers({
  token,
  mail,
  registration,
  loginStatus,
  recoverStatus,
  renewStatus
});
