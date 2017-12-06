import Constants from '../../constants';

export default {
  register: (username, password) => ({
		type: Constants.REGISTRATION_REQUEST,
		username,
		password,
	}),
  clearRegister: () => ({
	  type: Constants.REGISTRATION_CLEAR
  })
};
