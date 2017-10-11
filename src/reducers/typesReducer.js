import initialState from './initialState';
import * as actionTypes from '../constants/actionTypes';

// Handles regions related actions
export default function (state = initialState.types, action) {
  switch (action.type) {
    case actionTypes.SET_STORAGE_TYPES:
      return {...state, ...action.types};
    default:
      return state;
  }
}
