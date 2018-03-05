import Constants from '../../constants';

/* istanbul ignore next */
const defaultValue = (mode) => {
  switch (mode) {
    case "cb_pie":
    case "cb_bar":
      return "product";
    case "s3_chart":
      return "storage";
    case "s3_infos":
    default:
      return null;
  }
};

export default (state={}, action) => {
  let filters = Object.assign({}, state);
  switch (action.type) {
    case Constants.DASHBOARD_INSERT_ITEMS_FILTER:
      return action.filters;
    case Constants.DASHBOARD_ADD_ITEM:
      filters[action.id] = defaultValue(action.props.type);
      return filters;
    case Constants.DASHBOARD_SET_ITEM_FILTER:
      filters[action.id] = action.filter;
      return filters;
    case Constants.DASHBOARD_RESET_ITEMS_FILTER:
      Object.keys(filters).forEach((key) => {
        filters[key] = defaultValue("");
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
