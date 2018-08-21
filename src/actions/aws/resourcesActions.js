import Constants from '../../constants';

export default {
	selectAccount: (accountId) => ({
		type: Constants.AWS_RESOURCES_ACCOUNT_SELECTION,
		accountId
	}),
	get: {
		EC2: (accountId) => ({
      type: Constants.AWS_RESOURCES_GET_EC2,
      accountId
    }),
    RDS: (accountId) => ({
      type: Constants.AWS_RESOURCES_GET_RDS,
      accountId
    })
	},
	clear: {
		EC2: () => ({type: Constants.AWS_RESOURCES_GET_EC2_CLEAR}),
    RDS: () => ({type: Constants.AWS_RESOURCES_GET_RDS_CLEAR,})
	}
};
