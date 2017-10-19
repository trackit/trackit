import Constants from '../../constants';

export default (state=null, action) => {
  switch (action.type) {
    case Constants.GET_USER_TOKEN_SUCCESS:
      return action.token;
    case Constants.CLEAN_USER_TOKEN_SUCCESS:
      return null;
    default:
      return state;
  }
};
