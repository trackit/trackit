import Constants from '../../../constants';

export default (state={}, action) => {
  switch (action.type) {
    case Constants.AWS_SET_S3_DATES:
      return action.dates;
    case Constants.AWS_CLEAR_S3_DATES:
      return {};
    default:
      return state;
  }
};
