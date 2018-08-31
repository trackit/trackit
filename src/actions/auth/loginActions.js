import Constants from '../../constants';

export default {
	login: (username, password, awsToken) => ({
		type: Constants.LOGIN_REQUEST,
		username,
		password,
		awsToken,
	}),
};
