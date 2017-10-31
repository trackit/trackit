import { select } from 'redux-saga/effects';

const getTokenFromState = (state) => (state.auth.token);

export const getToken = () => {
  return select(getTokenFromState);
};
