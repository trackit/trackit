import Constants from '../../constants';

export default (state=[], action) => {
  switch (action.type) {
    case Constants.AWS_GET_ACCESS_SUCCESS:
      return action.access;
    default:
      return state;
  }
};
