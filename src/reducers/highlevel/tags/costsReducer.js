import Constants from '../../../constants';

export default (state={}, action) => {
  switch (action.type) {
    case Constants.HIGHLEVEL_TAGS_COST_REQUEST:
      return { status: false };
    case Constants.HIGHLEVEL_TAGS_COST_SUCCESS:
      return { status: true, values: action.values };
    case Constants.HIGHLEVEL_TAGS_COST_ERROR:
      return { status: true, error: action.error };
    case Constants.HIGHLEVEL_TAGS_COST_CLEAR:
      return { status: true, values: {} };
    default:
      return state;
  }
};
