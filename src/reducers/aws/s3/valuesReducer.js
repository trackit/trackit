import Constants from '../../../constants';

export default (state={}, action) => {
  switch (action.type) {
    case Constants.AWS_GET_S3_DATA:
      return { status: false };
    case Constants.AWS_GET_S3_DATA_SUCCESS:
      return { status: true, values: action.data };
    case Constants.AWS_GET_S3_DATA_ERROR:
      return { status: true, error: action.error };
    default:
      return state;
  }
};
