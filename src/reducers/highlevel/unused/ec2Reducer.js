import Constants from '../../../constants';

export default (state={}, action) => {
  switch (action.type) {
    case Constants.HIGHLEVEL_UNUSED_EC2_REQUEST:
      return { status: false };
    case Constants.HIGHLEVEL_UNUSED_EC2_SUCCESS:
      return { status: true, values: action.data };
    case Constants.HIGHLEVEL_UNUSED_EC2_ERROR:
      return { status: true, error: action.error };
    default:
      return state;
  }
};
