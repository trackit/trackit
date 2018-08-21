import { combineReducers } from 'redux';
import account from './accountReducer';
import EC2 from './EC2Reducer';
import RDS from './RDSReducer';

export default combineReducers({
  account,
  EC2,
  RDS
});
