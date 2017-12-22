import Constants from '../../../constants';

export default (state=null, action) => {
  switch (action.type) {
    case Constants.AWS_CLEAR_COSTS_DATES:
      return null;
    case Constants.AWS_SET_COSTS_DATES:
    return action.dates;
    default:
      return state;
  }
};
