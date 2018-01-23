import Constants from '../../../constants';

const defaultValue = "product";

export default (state={}, action) => {
  let filters = Object.assign({}, state);
  switch (action.type) {
    case Constants.AWS_CLEAR_COSTS_FILTER:
      return {};
    case Constants.AWS_RESET_COSTS_FILTER:
      Object.keys(filters).forEach((key) => {
        filters[key] = defaultValue;
      });
      return filters;
    case Constants.AWS_SET_COSTS_FILTER:
      filters[action.id] = action.filter;
      return filters;
    default:
      return state;
  }
};
