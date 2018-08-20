import { combineReducers } from 'redux';
import all from './allReducer';
import selection from './selectionReducer';
import external from './externalReducer';
import bills from './billsReducer';
import creation from './creationReducer';
import billCreation from './billCreationReducer';
import billsStatus from './billsStatusReducer';
import billEdition from './billEditionReducer';

export default combineReducers({
  all,
  selection,
  external,
  bills,
  creation,
  billCreation,
  billsStatus,
  billEdition
});
