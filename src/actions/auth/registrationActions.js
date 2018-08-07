import Constants from '../../constants';

export default {
  register: (username, password, awsToken) => ({
		type: Constants.REGISTRATION_REQUEST,
		username,
		password,
		awsToken
	}),
  clearRegister: () => ({
	  type: Constants.REGISTRATION_CLEAR
  })
};
