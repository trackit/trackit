import Constants from '../../constants';

export default {
	login: (username, password) => ({
		type: Constants.LOGIN_REQUEST,
		username,
		password,
	}),
};
