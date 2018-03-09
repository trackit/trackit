import Constants from '../../constants';

const defaultValue = "day";
const defaultValuePie = "month";

export default (state={}, action) => {
  let intervals = Object.assign({}, state);
  switch (action.type) {
    case Constants.DASHBOARD_INSERT_ITEMS_INTERVAL:
      return action.intervals;
    case Constants.DASHBOARD_ADD_ITEM:
      switch (action.chartType) {
        case "pie":
          intervals[action.id] = defaultValuePie;
          break;
        case "bar":
        default:
          intervals[action.id] = defaultValue;
      }
      return intervals;
    case Constants.DASHBOARD_SET_ITEM_INTERVAL:
      intervals[action.id] = action.interval;
      return intervals;
    case Constants.DASHBOARD_RESET_ITEMS_INTERVAL:
      Object.keys(intervals).forEach((key) => {
        intervals[key] = defaultValue;
      });
      return intervals;
    case Constants.DASHBOARD_REMOVE_ITEM:
      if (intervals.hasOwnProperty(action.id))
        delete intervals[action.id];
      return intervals;
    case Constants.DASHBOARD_CLEAR_ITEMS_INTERVAL:
      return {};
    default:
      return state;
  }
};
