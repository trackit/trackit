import Constants from '../../constants';

export default {
  initCharts: () => ({type: Constants.AWS_INIT_CHARTS}),
  addChart: (id, chartType) => ({type: Constants.AWS_ADD_CHART, id, chartType}),
  removeChart: (id) => ({type: Constants.AWS_REMOVE_CHART, id}),
	getCosts: (id, begin, end, filters) => ({
		type: Constants.AWS_GET_COSTS,
    id,
    begin,
    end,
    filters
	}),
  setCostsDates: (id, startDate, endDate) => ({
    type: Constants.AWS_SET_COSTS_DATES,
    id,
    dates: {
      startDate,
      endDate
    }
  }),
  clearCostsDates: () => ({type: Constants.AWS_CLEAR_COSTS_DATES}),
  resetCostsDates: () => ({type: Constants.AWS_RESET_COSTS_DATES}),
  setCostsInterval: (id, interval) => ({
    type: Constants.AWS_SET_COSTS_INTERVAL,
    id,
    interval
  }),
  clearCostsInterval: () => ({type: Constants.AWS_CLEAR_COSTS_INTERVAL}),
  resetCostsInterval: () => ({type: Constants.AWS_RESET_COSTS_INTERVAL}),
  setCostsFilter: (id, filter) => ({
    type: Constants.AWS_SET_COSTS_FILTER,
    id,
    filter
  }),
  clearCostsFilter: () => ({type: Constants.AWS_CLEAR_COSTS_FILTER}),
  resetCostsFilter: () => ({type: Constants.AWS_RESET_COSTS_FILTER}),
};
