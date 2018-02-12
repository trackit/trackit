import { combineReducers } from 'redux';
import all from './allReducer';
import selection from './selectionReducer';
import external from './externalReducer';
import bills from './billsReducer';
import creation from './creationReducer';
import billCreation from './billCreationReducer';
import retrieved from './retrievedReducer';

export default combineReducers({
  all,
  selection,
  external,
  bills,
  retrieved,
  creation,
  billCreation
});
