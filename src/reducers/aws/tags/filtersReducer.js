import Constants from '../../../constants';

const defaultValue = "product";

export default (state={}, action) => {
  let filters = Object.assign({}, state);
  switch (action.type) {
    case Constants.AWS_TAGS_INSERT_FILTERS:
      return action.filters;
    case Constants.AWS_TAGS_ADD_CHART:
      filters[action.id] = defaultValue;
      return filters;
    case Constants.AWS_TAGS_SET_FILTER:
      filters[action.id] = action.filter;
      return filters;
    case Constants.AWS_TAGS_RESET_FILTERS:
      Object.keys(filters).forEach((key) => {
        filters[key] = defaultValue;
      });
      return filters;
    case Constants.AWS_TAGS_REMOVE_CHART:
      if (filters.hasOwnProperty(action.id))
        delete filters[action.id];
      return filters;
    case Constants.AWS_TAGS_CLEAR_FILTERS:
      return {};
    default:
      return state;
  }
};
