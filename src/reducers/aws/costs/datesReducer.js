import Constants from '../../../constants';
import moment from "moment/moment";

const defaultValue = {
  startDate: moment().subtract(1, 'month').startOf('month'),
  endDate: moment().subtract(1, 'month').endOf('month')
};

export default (state={}, action) => {
  let dates = Object.assign({}, state);
  switch (action.type) {
    case Constants.AWS_CLEAR_COSTS_DATES:
      return {};
    case Constants.AWS_RESET_COSTS_DATES:
      Object.keys(dates).forEach((key) => {
        dates[key] = defaultValue;
      });
      return dates;
    case Constants.AWS_SET_COSTS_DATES:
      dates[action.id] = action.dates;
      return dates;
    default:
      return state;
  }
};
