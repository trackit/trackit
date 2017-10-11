import { fork } from 'redux-saga/effects';
import * as Watchers from './watcher';

// Here, we register our watcher saga(s) and export as a single generator
// function (startForeman) as our root Saga.
export default function* startForman() {
  for (var key in Watchers) {
    yield fork(Watchers[key]);
  }
}
