import Constants from '../../../constants';

const defaultValue = {status: true, value: null};

export default (state=defaultValue, action) => {
  switch (action.type) {
    case Constants.AWS_RESOURCES_GET_ES_CLEAR:
      return defaultValue;
    case Constants.AWS_RESOURCES_GET_ES:
      return {status: false};
    case Constants.AWS_RESOURCES_GET_ES_SUCCESS:
      return {status: true, value: action.report};
    case Constants.AWS_RESOURCES_GET_ES_ERROR:
      return {status: true, error: action.error};
    default:
      return state;
  }
};
