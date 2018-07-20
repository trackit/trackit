import Constants from '../../constants';

export default {
  recover: (username) => ({
    type: Constants.RECOVER_PASSWORD_REQUEST,
    username
  }),
  clearRecover: () => ({
    type: Constants.RECOVER_PASSWORD_CLEAR
  }),
  renew: (id, password, token) => ({
    type: Constants.RENEW_PASSWORD_REQUEST,
    id,
	  password,
	  token
  }),
  clearRenew: () => ({
    type: Constants.RENEW_PASSWORD_CLEAR
  }),
};
