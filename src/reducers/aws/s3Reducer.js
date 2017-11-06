import Constants from '../../constants';

export default (state=[], action) => {
  switch (action.type) {
    case Constants.AWS_GET_S3_DATA_SUCCESS:
      return action.s3Data;
    case Constants.AWS_GET_S3_DATA_ERROR:
      return [];
    default:
      return state;
  }
};
