import Constants from '../../constants';

export default {
	setDates: (startDate, endDate) => ({
		type: Constants.EVENTS_SET_DATES,
		dates: {
			startDate,
            endDate
		}
	}),
	clearDates: () => ({type: Constants.EVENTS_CLEAR_DATES}),
	getData: (begin, end) => ({
		type: Constants.GET_EVENTS_DATA,
		begin,
		end
	}),
};
