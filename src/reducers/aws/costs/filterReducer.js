import Constants from '../../../constants';

export default (state={}, action) => {
  switch (action.type) {
    case Constants.AWS_CLEAR_COSTS_FILTER:
      return {};
    case Constants.AWS_SET_COSTS_FILTER:
      let filters = Object.assign({}, state);
      filters[action.id] = action.filter;
    return filters;
    default:
      return state;
  }
};
