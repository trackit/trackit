import { combineReducers } from 'redux';
import account from './accountReducer';
import reportList from './reportListReducer';
import download from './downloadReducer';

export default combineReducers({
  account,
  reportList,
  download
});
