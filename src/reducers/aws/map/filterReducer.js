import Constants from '../../../constants';

const defaultValue = "region";

export default (state=defaultValue, action) => {
  switch (action.type) {
    case Constants.AWS_MAP_SET_FILTER:
      return action.filter;
    case Constants.AWS_MAP_RESET_FILTER:
      return defaultValue;
    default:
      return state;
  }
};
