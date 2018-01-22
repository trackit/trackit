import Constants from '../../../constants';

export default (state={}, action) => {
  switch (action.type) {
    case Constants.AWS_CLEAR_COSTS_INTERVAL:
      return {};
    case Constants.AWS_SET_COSTS_INTERVAL:
      let intervals = Object.assign({}, state);
      intervals[action.id] = action.interval;
      return intervals;
    default:
      return state;
  }
};
