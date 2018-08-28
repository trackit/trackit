import Constants from '../../../constants';

const defaultValue = {status: true, value: null};

export default (state=defaultValue, action) => {
  switch (action.type) {
    case Constants.AWS_RESOURCES_GET_EC2_CLEAR:
      return defaultValue;
    case Constants.AWS_RESOURCES_GET_EC2:
      return {status: false};
    case Constants.AWS_RESOURCES_GET_EC2_SUCCESS:
      return {status: true, value: action.report};
    case Constants.AWS_RESOURCES_GET_EC2_ERROR:
      return {status: true, error: action.error};
    default:
      return state;
  }
};
