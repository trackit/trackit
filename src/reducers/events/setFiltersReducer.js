import Constants from '../../constants';

export default (state={}, action) => {
  switch (action.type) {
    case Constants.EVENTS_SET_FILTERS:
      return { status: false };
    case Constants.EVENTS_SET_FILTERS_CLEAR:
      return { status: true, values: null };
    case Constants.EVENTS_SET_FILTERS_SUCCESS:
      return { status: true, values: action.data };
    case Constants.EVENTS_SET_FILTERS_ERROR:
      return { status: true, error: action.error };
    default:
      return state;
  }
};
