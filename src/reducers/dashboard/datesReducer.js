import Constants from '../../constants';
import moment from "moment/moment";

const defaultValue = {
  startDate: moment().subtract(1, 'month').startOf('month'),
  endDate: moment().subtract(1, 'month').endOf('month')
};

export default (state={}, action) => {
  switch (action.type) {
    case Constants.DASHBOARD_RESET_DATES:
      return defaultValue;
    case Constants.DASHBOARD_INSERT_DATES:
      return {
        startDate: moment(action.dates.startDate),
        endDate: moment(action.dates.endDate),
      };
    case Constants.DASHBOARD_SET_DATES:
      return action.dates;
    case Constants.DASHBOARD_CLEAR_DATES:
      return {};
    default:
      return state;
  }
};
