import Constants from '../../../constants';

export default (state={}, action) => {
  switch (action.type) {
    case Constants.HIGHLEVEL_TAGS_KEYS_REQUEST:
      return { status: false };
    case Constants.HIGHLEVEL_TAGS_KEYS_SUCCESS:
      return { status: true, values: action.keys };
    case Constants.HIGHLEVEL_TAGS_KEYS_ERROR:
      return { status: true, error: action.error };
    default:
      return state;
  }
};
