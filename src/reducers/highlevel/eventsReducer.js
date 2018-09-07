import Constants from '../../constants';

export default (state={}, action) => {
  switch (action.type) {
    case Constants.HIGHLEVEL_EVENTS_REQUEST:
      return { status: false };
    case Constants.HIGHLEVEL_EVENTS_SUCCESS:
      return { status: true, values: action.events };
    case Constants.HIGHLEVEL_EVENTS_ERROR:
      return { status: true, error: action.error };
    default:
      return state;
  }
};
