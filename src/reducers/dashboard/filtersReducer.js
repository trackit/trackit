import Constants from '../../constants';

const defaultValue = "product";

export default (state={}, action) => {
  let filters = Object.assign({}, state);
  switch (action.type) {
    case Constants.DASHBOARD_INSERT_ITEMS_FILTER:
      return action.filters;
    case Constants.DASHBOARD_ADD_ITEM:
      filters[action.id] = defaultValue;
      return filters;
    case Constants.DASHBOARD_SET_ITEM_FILTER:
      filters[action.id] = action.filter;
      return filters;
    case Constants.DASHBOARD_RESET_ITEMS_FILTER:
      Object.keys(filters).forEach((key) => {
        filters[key] = defaultValue;
      });
      return filters;
    case Constants.DASHBOARD_REMOVE_ITEM:
      if (filters.hasOwnProperty(action.id))
        delete filters[action.id];
      return filters;
    case Constants.DASHBOARD_CLEAR_ITEMS_FILTER:
      return {};
    default:
      return state;
  }
};
