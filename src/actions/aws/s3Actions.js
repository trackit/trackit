import Constants from '../../constants';

export default {
	getS3Data: () => ({
		type: Constants.AWS_GET_S3_DATA
	}),
	setS3ViewDates: (startDate, endDate) => ({
		type: Constants.AWS_SET_S3_VIEW_DATES,
		startDate,
		endDate,
	}),
};
