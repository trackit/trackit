import { combineReducers } from 'redux';
import dates from './datesReducer';
import EC2 from './EC2Reducer';
import RDS from './RDSReducer';
import ES from './ESReducer';

export default combineReducers({
  dates,
  EC2,
  RDS,
  ES
});
