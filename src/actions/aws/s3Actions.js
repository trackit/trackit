import Constants from '../../constants';

export default {
	getData: (begin, end) => ({
		type: Constants.AWS_GET_S3_DATA,
		begin,
		end
	}),
	setDates: (startDate, endDate) => ({
		type: Constants.AWS_SET_S3_DATES,
		dates: {
			startDate,
      endDate
		}
	}),
	clearDates: () => ({type: Constants.AWS_CLEAR_S3_DATES})
};
