import initialState from './initialState';
import * as actionTypes from '../constants/actionTypes';

// Handles regions related actions
export default function (state = initialState.aws, action) {
  switch (action.type) {
    case actionTypes.GET_PRICING_AWS_SUCCESS:
      return {...state, pricing: action.pricing};
    default:
      return state;
  }
}
