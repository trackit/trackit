import { createStore, applyMiddleware, compose } from 'redux';
import createSagaMiddleware from 'redux-saga';
import RavenMiddleware from 'redux-raven-middleware';

import Config from '../config';

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

  const middlewares = (Config.sentryDSN ? applyMiddleware(sagaMiddleware, RavenMiddleware(Config.sentryDSN)) : applyMiddleware(sagaMiddleware));

  return {
    ...createStore(
      rootReducer,
      initialState,
      composeEnhancers(middlewares)
    ),
    runSaga: sagaMiddleware.run(rootSaga)
  };
};
