import Constants from '../../constants';

export default {
  initCharts: () => ({type: Constants.AWS_TAGS_INIT_CHARTS}),
  addChart: (id) => ({type: Constants.AWS_TAGS_ADD_CHART, id}),
  removeChart: (id) => ({type: Constants.AWS_TAGS_REMOVE_CHART, id}),
  getValues: (id, begin, end, filter, key) => ({
    type: Constants.AWS_TAGS_GET_VALUES,
    id,
    begin,
    end,
    filter,
    key
  }),
  getKeys: (id, begin, end) => ({
    type: Constants.AWS_TAGS_GET_KEYS,
    id,
    begin,
    end
  }),
  clearKeys: () => ({type: Constants.AWS_TAGS_GET_KEYS_CLEAR}),
  selectKey: (id, key) => ({type: Constants.AWS_TAGS_SELECT_KEY, id, key}),
  setDates: (id, startDate, endDate) => ({
    type: Constants.AWS_TAGS_SET_DATES,
    id,
    dates: {
      startDate,
      endDate
    }
  }),
  clearDates: () => ({type: Constants.AWS_TAGS_CLEAR_DATES}),
  resetDates: () => ({type: Constants.AWS_TAGS_RESET_DATES}),
  selectFilter: (id, filter) => ({
    type: Constants.AWS_TAGS_SET_FILTER,
    id,
    filter
  }),
  clearFilters: () => ({type: Constants.AWS_TAGS_CLEAR_FILTERS}),
  resetFilters: () => ({type: Constants.AWS_TAGS_RESET_FILTERS}),
};
