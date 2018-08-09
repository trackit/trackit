import Constants from '../../../constants';

const defaultValue = "month";

export default (state={}, action) => {
  let intervals = Object.assign({}, state);
  switch (action.type) {
    case Constants.AWS_TAGS_INSERT_INTERVAL:
      return action.interval;
    case Constants.AWS_TAGS_ADD_CHART:
      intervals[action.id] = defaultValue;
      return intervals;
    case Constants.AWS_TAGS_SET_INTERVAL:
      intervals[action.id] = action.interval;
      return intervals;
    case Constants.AWS_TAGS_RESET_INTERVAL:
      Object.keys(intervals).forEach((key) => {
        intervals[key] = defaultValue;
      });
      return intervals;
    case Constants.AWS_TAGS_REMOVE_CHART:
      if (intervals.hasOwnProperty(action.id))
        delete intervals[action.id];
      return intervals;
    case Constants.AWS_TAGS_CLEAR_INTERVAL:
      return {};
    default:
      return state;
  }
};
