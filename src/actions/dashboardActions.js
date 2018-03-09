import Constants from '../constants';

export default {
  initDashboard: () => ({type: Constants.DASHBOARD_INIT_ITEMS}),
  updateDashboard: (items) => ({type: Constants.DASHBOARD_UPDATE_ITEMS, items}),
  addItem: (id, props) => ({type: Constants.DASHBOARD_ADD_ITEM, id, props}),
  removeItem: (id) => ({type: Constants.DASHBOARD_REMOVE_ITEM, id}),
	getData: (id, itemType, begin, end, filters) => ({
		type: Constants.DASHBOARD_GET_VALUES,
    id,
    itemType,
    begin,
    end,
    filters
	}),
  setItemDates: (id, startDate, endDate) => ({
    type: Constants.DASHBOARD_SET_ITEM_DATES,
    id,
    dates: {
      startDate,
      endDate
    }
  }),
  clearItemDates: () => ({type: Constants.DASHBOARD_CLEAR_ITEMS_DATES}),
  resetItemDates: () => ({type: Constants.DASHBOARD_RESET_ITEMS_DATES}),
  setItemInterval: (id, interval) => ({
    type: Constants.DASHBOARD_SET_ITEM_INTERVAL,
    id,
    interval
  }),
  clearItemInterval: () => ({type: Constants.DASHBOARD_CLEAR_ITEMS_INTERVAL}),
  resetItemInterval: () => ({type: Constants.DASHBOARD_RESET_ITEMS_INTERVAL}),
  setItemFilter: (id, filter) => ({
    type: Constants.DASHBOARD_SET_ITEM_FILTER,
    id,
    filter
  }),
  clearItemFilter: () => ({type: Constants.DASHBOARD_CLEAR_ITEMS_FILTER}),
  resetItemFilter: () => ({type: Constants.DASHBOARD_RESET_ITEMS_FILTER}),
};
