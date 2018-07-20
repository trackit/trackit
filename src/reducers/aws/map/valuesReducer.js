import Constants from '../../../constants';

export default (state={}, action) => {
  switch (action.type) {
    case Constants.AWS_MAP_GET_COSTS_ERROR:
      return { status: true, error: action.error };
    case Constants.AWS_MAP_GET_COSTS:
      return { status: false };
    case Constants.AWS_MAP_GET_COSTS_SUCCESS:
      return { status: true, values: action.costs };
    case Constants.AWS_MAP_GET_COSTS_CLEAR:
      return { status: true, values: null };
    default:
      return state;
  }
};
