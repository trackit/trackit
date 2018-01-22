import Constants from '../../constants';

export default {
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
  setCostsInterval: (id, interval) => ({
    type: Constants.AWS_SET_COSTS_INTERVAL,
    id,
    interval
  }),
  clearCostsInterval: () => ({type: Constants.AWS_CLEAR_COSTS_INTERVAL}),
  setCostsFilter: (id, filter) => ({
    type: Constants.AWS_SET_COSTS_FILTER,
    id,
    filter
  }),
  clearCostsFilter: () => ({type: Constants.AWS_CLEAR_COSTS_FILTER}),
};
