import Constants from '../../constants';
import moment from "moment/moment";

const defaultValue = {
  startDate: moment().subtract(1, 'month').startOf('month'),
  endDate: moment().subtract(1, 'month').endOf('month')
};

export default (state={}, action) => {
  let dates = Object.assign({}, state);
  switch (action.type) {
    case Constants.DASHBOARD_INSERT_ITEMS_DATES:
      let newDates = Object.assign({}, action.dates);
      Object.keys(newDates).forEach((id) => {
        newDates[id].startDate = moment(newDates[id].startDate);
        newDates[id].endDate = moment(newDates[id].endDate);
      });
      return newDates;
    case Constants.DASHBOARD_ADD_ITEM:
      dates[action.id] = defaultValue;
      return dates;
    case Constants.DASHBOARD_SET_ITEM_DATES:
      dates[action.id] = action.dates;
      return dates;
    case Constants.DASHBOARD_RESET_ITEMS_DATES:
      Object.keys(dates).forEach((key) => {
        dates[key] = defaultValue;
      });
      return dates;
    case Constants.DASHBOARD_REMOVE_ITEM:
      if (dates.hasOwnProperty(action.id))
        delete dates[action.id];
      return dates;
    case Constants.DASHBOARD_CLEAR_ITEMS_DATES:
      return {};
    default:
      return state;
  }
};
