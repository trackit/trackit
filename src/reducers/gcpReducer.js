import initialState from './initialState';
import * as actionTypes from '../constants/actionTypes';

export default function (state = initialState.gcp, action) {
  switch (action.type) {
    case actionTypes.GET_PRICING_GCP_SUCCESS:
      return {...state, pricing: action.pricing};
    default:
      return state;
  }
}
