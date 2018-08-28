import Constants from '../../constants';

export default {
	get: {
		EC2: () => ({type: Constants.AWS_RESOURCES_GET_EC2}),
    RDS: () => ({type: Constants.AWS_RESOURCES_GET_RDS})
	},
	clear: {
		EC2: () => ({type: Constants.AWS_RESOURCES_GET_EC2_CLEAR}),
    RDS: () => ({type: Constants.AWS_RESOURCES_GET_RDS_CLEAR,})
	}
};
