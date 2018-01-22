import Constants from '../../../constants';

export default (state={}, action) => {
  switch (action.type) {
    case Constants.AWS_GET_COSTS_ERROR:
      return {};
    case Constants.AWS_GET_COSTS_SUCCESS:
      let costs = Object.assign({}, state);
      costs[action.id] = action.costs;
      return costs;
    default:
      return state;
  }
};
