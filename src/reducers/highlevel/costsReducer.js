import Constants from '../../constants';

export default (state={}, action) => {
  switch (action.type) {
    case Constants.HIGHLEVEL_COSTS_REQUEST:
      return { status: false };
    case Constants.HIGHLEVEL_COSTS_SUCCESS:
      return { status: true, values: { months: action.months, history: action.history } };
    case Constants.HIGHLEVEL_COSTS_ERROR:
      return { status: true, error: action.error };
    default:
      return state;
  }
};
