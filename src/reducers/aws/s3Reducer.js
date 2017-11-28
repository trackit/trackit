import Constants from '../../constants';

export default (state={}, action) => {
  switch (action.type) {
    case Constants.AWS_GET_S3_DATA_SUCCESS:
      return {
        ...state,
        data: action.s3Data
      };
    case Constants.AWS_GET_S3_DATA_ERROR:
      return {
        ...state,
        data: []
      };
    case Constants.AWS_SET_S3_VIEW_DATES:
    return {
      ...state,
      view: {
        startDate: action.startDate,
        endDate: action.endDate
      }
    };
    default:
      return state;
  }
};
