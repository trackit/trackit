import Constants from '../../../constants';

export default (state=null, action) => {
  switch (action.type) {
    case Constants.HIGHLEVEL_TAGS_KEYS_SELECT:
      return action.key;
    case Constants.HIGHLEVEL_TAGS_KEYS_CLEAR_SELECTED:
      return null;
    default:
      return state;
  }
};
