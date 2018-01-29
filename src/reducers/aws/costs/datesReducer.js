import Constants from '../../../constants';
import moment from "moment/moment";

const defaultValue = {
  startDate: moment().subtract(1, 'month').startOf('month'),
  endDate: moment().subtract(1, 'month').endOf('month')
};

export default (state={}, action) => {
  let dates = Object.assign({}, state);
  switch (action.type) {
    case Constants.AWS_INSERT_COSTS_DATES:
      let newDates = Object.assign({}, action.dates);
      Object.keys(newDates).forEach((id) => {
        newDates[id].startDate = moment(newDates[id].startDate);
        newDates[id].endDate = moment(newDates[id].endDate);
      });
      return newDates;
    case Constants.AWS_ADD_CHART:
      dates[action.id] = defaultValue;
      return dates;
    case Constants.AWS_SET_COSTS_DATES:
      dates[action.id] = action.dates;
      return dates;
    case Constants.AWS_RESET_COSTS_DATES:
      Object.keys(dates).forEach((key) => {
        dates[key] = defaultValue;
      });
      return dates;
    case Constants.AWS_REMOVE_CHART:
      if (dates.hasOwnProperty(action.id))
        delete dates[action.id];
      return dates;
    case Constants.AWS_CLEAR_COSTS_DATES:
      return {};
    default:
      return state;
  }
};
