import Constants from '../../constants';

export default {
	selectAccount: (accountId) => ({
		type: Constants.AWS_REPORTS_ACCOUNT_SELECTION,
		accountId
	}),
	requestGetReports: (accountId) => ({
		type: Constants.AWS_GET_REPORTS_REQUESTED,
		accountId
	}),
	requestDownloadReport: (accountId, reportType, fileName) => ({
		type: Constants.AWS_DOWNLOAD_REPORT_REQUESTED,
		accountId,
		reportType,
		fileName
	}),
};
