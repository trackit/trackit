import Constants from '../../constants';

export default (state={}, action) => {
  switch (action.type) {
    case Constants.EVENTS_SET_DATES:
      return action.dates;
    case Constants.EVENTS_CLEAR_DATES:
      return {};
    default:
      return state;
  }
};
