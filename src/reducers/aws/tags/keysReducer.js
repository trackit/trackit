import Constants from '../../../constants';

export default (state={}, action) => {
  let keys = Object.assign({}, state);
  switch (action.type) {
    case Constants.AWS_TAGS_GET_KEYS_ERROR:
      keys[action.id] = { status: true, error: action.error };
      return keys;
    case Constants.AWS_TAGS_GET_KEYS:
      keys[action.id] = { status: false };
      return keys;
    case Constants.AWS_TAGS_GET_KEYS_SUCCESS:
      keys[action.id] = { status: true, values: action.tags };
      return keys;
    case Constants.AWS_TAGS_GET_KEYS_CLEAR:
      keys[action.id] = { status: true, values: null };
      return keys;
    default:
      return state;
  }
};
