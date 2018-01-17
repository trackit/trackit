import Constants from '../../../constants';

export default (state=null, action) => {
  switch (action.type) {
    case Constants.AWS_CLEAR_COSTS_INTERVAL:
      return null;
    case Constants.AWS_SET_COSTS_INTERVAL:
    return action.interval;
    default:
      return state;
  }
};
