import Constants from '../../constants';

export default {
  initCharts: () => ({type: Constants.AWS_TAGS_INIT_CHARTS}),
  addChart: (id) => ({type: Constants.AWS_TAGS_ADD_CHART, id}),
  removeChart: (id) => ({type: Constants.AWS_TAGS_REMOVE_CHART, id}),
  getValues: (id, begin, end, key) => ({
    type: Constants.AWS_TAGS_GET_VALUES,
    id,
    begin,
    end,
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
  setInterval: (id, interval) => ({
    type: Constants.AWS_TAGS_SET_INTERVAL,
    id,
    interval
  }),
  clearInterval: () => ({type: Constants.AWS_TAGS_CLEAR_INTERVAL}),
  resetInterval: () => ({type: Constants.AWS_TAGS_RESET_INTERVAL}),
};
