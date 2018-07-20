import Constants from '../../constants';

export default {
	getCosts: (begin, end) => ({
		type: Constants.AWS_MAP_GET_COSTS,
    begin,
    end
	}),
  clearCosts: () => ({type: Constants.AWS_MAP_GET_COSTS_CLEAR}),
  setDates: (startDate, endDate) => ({
    type: Constants.AWS_MAP_SET_COSTS_DATES,
    dates: {
      startDate,
      endDate
    }
  }),
  clearDates: () => ({type: Constants.AWS_MAP_CLEAR_COSTS_DATES}),
  resetDates: () => ({type: Constants.AWS_MAP_RESET_COSTS_DATES}),
};
