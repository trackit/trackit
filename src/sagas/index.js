import { fork } from 'redux-saga/effects';
import Watchers from './watcher';

export default function* startForman() {
  for (let key in Watchers) {
    yield fork(Watchers[key]);
  }
};
