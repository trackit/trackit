import Constants from '../../constants';

export default {
	token: {
    get: () => ({
      type: Constants.GET_USER_TOKEN
    }),
    clean: () => ({
      type: Constants.CLEAN_USER_TOKEN
    })
  }
};
