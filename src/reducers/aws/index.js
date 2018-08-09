import { combineReducers } from 'redux';
import accounts from './accounts';
import costs from './costs';
import s3 from './s3';
import reports from './reports'
import resources from './resources'
import map from './map';
import tags from './tags';

export default combineReducers({
  accounts,
  s3,
  costs,
  reports,
  map,
  resources,
  tags
});
