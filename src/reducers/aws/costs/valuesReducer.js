import Constants from '../../../constants';

export default (state=null, action) => {
  switch (action.type) {
    case Constants.AWS_GET_COSTS:
    case Constants.AWS_GET_COSTS_ERROR:
      return null;
    case Constants.AWS_GET_COSTS_SUCCESS:
    return action.costs;
    default:
      return state;
  }
};
