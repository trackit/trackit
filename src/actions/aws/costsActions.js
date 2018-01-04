import Constants from '../../constants';

export default {
	getCosts: (begin, end, filters) => ({
		type: Constants.AWS_GET_COSTS,
    begin,
    end,
    filters
	}),
  setCostsDates: (startDate, endDate) => ({
    type: Constants.AWS_SET_COSTS_DATES,
    dates: {
      startDate,
      endDate
    }
  }),
  clearCostsDates: () => ({type: Constants.AWS_CLEAR_COSTS_DATES}),
  setCostsInterval: (interval) => ({
    type: Constants.AWS_SET_COSTS_INTERVAL,
    interval
  }),
  clearCostsInterval: () => ({type: Constants.AWS_CLEAR_COSTS_INTERVAL}),
  setCostsFilter: (filter) => ({
    type: Constants.AWS_SET_COSTS_FILTER,
    filter
  }),
  clearCostsFilter: () => ({type: Constants.AWS_CLEAR_COSTS_FILTER}),
};
