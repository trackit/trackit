import Constants from '../../constants';

export default {
	getCosts: (begin, end, filters, accounts=undefined) => ({
		type: Constants.AWS_GET_COSTS,
    begin,
    end,
    filters,
    accounts
	}),
  setCostsDates: (startDate, endDate) => ({
    type: Constants.AWS_SET_COSTS_DATES,
    dates: {
      startDate,
      endDate
    }
  }),
  setCostsInterval: (interval) => ({
    type: Constants.AWS_SET_COSTS_INTERVAL,
    interval
  })
};
