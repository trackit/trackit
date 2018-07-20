import Constants from '../../../constants';
import moment from "moment/moment";

const defaultValue = {
  startDate: moment().subtract(1, 'month').startOf('month'),
  endDate: moment().subtract(1, 'month').endOf('month')
};

export default (state={}, action) => {
  switch (action.type) {
    case Constants.AWS_MAP_SET_COSTS_DATES:
      return action.dates;
    case Constants.AWS_MAP_RESET_COSTS_DATES:
      return defaultValue;
    case Constants.AWS_MAP_CLEAR_COSTS_DATES:
      return {};
    default:
      return state;
  }
};
