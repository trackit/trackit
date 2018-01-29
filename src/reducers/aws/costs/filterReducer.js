import Constants from '../../../constants';

const defaultValue = "product";

export default (state={}, action) => {
  let filters = Object.assign({}, state);
  switch (action.type) {
    case Constants.AWS_INSERT_COSTS_FILTER:
      return action.filter;
    case Constants.AWS_ADD_CHART:
      filters[action.id] = defaultValue;
      return filters;
    case Constants.AWS_SET_COSTS_FILTER:
      filters[action.id] = action.filter;
      return filters;
    case Constants.AWS_RESET_COSTS_FILTER:
      Object.keys(filters).forEach((key) => {
        filters[key] = defaultValue;
      });
      return filters;
    case Constants.AWS_REMOVE_CHART:
      if (filters.hasOwnProperty(action.id))
        delete filters[action.id];
      return filters;
    case Constants.AWS_CLEAR_COSTS_FILTER:
      return {};
    default:
      return state;
  }
};
