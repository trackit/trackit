import Constants from '../../constants';

export default (state=[], action) => {
  switch (action.type) {
    case Constants.GET_USER_TOKEN_SUCCESS:
      return action.token;
    default:
      return state;
  }
};
