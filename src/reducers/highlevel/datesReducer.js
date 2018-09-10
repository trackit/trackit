import Constants from '../../constants';

export default (state={}, action) => {
  switch (action.type) {
    case Constants.HIGHLEVEL_SET_DATES:
      return action.dates;
    case Constants.HIGHLEVEL_CLEAR_DATES:
      return {};
    default:
      return state;
  }
};
