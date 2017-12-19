import Constants from '../../../constants';

export default (state=null, action) => {
  switch (action.type) {
    case Constants.AWS_CLEAR_COSTS_FILTER:
      return null;
    case Constants.AWS_SET_COSTS_FILTER:
    return action.filter;
    default:
      return state;
  }
};
