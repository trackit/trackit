import Constants from '../../../constants';

export default (state={}, action) => {
  switch (action.type) {
    case Constants.AWS_CLEAR_COSTS_DATES:
      return {};
    case Constants.AWS_SET_COSTS_DATES:
      let dates = Object.assign({}, state);
      dates[action.id] = action.dates;
      return dates;
    default:
      return state;
  }
};
