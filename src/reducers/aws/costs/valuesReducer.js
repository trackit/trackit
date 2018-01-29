import Constants from '../../../constants';

export default (state={}, action) => {
  let costs = Object.assign({}, state);
  switch (action.type) {
    case Constants.AWS_GET_COSTS_ERROR:
      return {};
    case Constants.AWS_GET_COSTS_SUCCESS:
      costs[action.id] = action.costs;
      return costs;
    case Constants.AWS_REMOVE_CHART:
      if (costs.hasOwnProperty(action.id))
        delete costs[action.id];
      return costs;
    default:
      return state;
  }
};
