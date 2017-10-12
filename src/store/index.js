import { createStore, applyMiddleware, compose } from 'redux';
import createSagaMiddleware from 'redux-saga';

import rootReducer from '../reducers';
import rootSaga from '../sagas';

import initialState from './initialState';

//  Returns the store instance
// It can  also take initialState argument when provided
export default () => {
  // Creating the redux-saga middleware
  const sagaMiddleware = createSagaMiddleware();

  // Enabling Redux Devtools
  const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

  return {
    ...createStore(
      rootReducer,
      initialState,
      composeEnhancers(
        applyMiddleware(sagaMiddleware)
      )
    ),
    runSaga: sagaMiddleware.run(rootSaga)
  };
};
