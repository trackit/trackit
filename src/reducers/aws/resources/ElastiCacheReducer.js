import Constants from '../../../constants';

const defaultValue = {status: true, value: null};

export default (state=defaultValue, action) => {
  switch (action.type) {
    case Constants.AWS_RESOURCES_GET_ELASTICACHE_CLEAR:
      return defaultValue;
    case Constants.AWS_RESOURCES_GET_ELASTICACHE:
      return {status: false};
    case Constants.AWS_RESOURCES_GET_ELASTICACHE_SUCCESS:
      return {status: true, value: action.report};
    case Constants.AWS_RESOURCES_GET_ELASTICACHE_ERROR:
      return {status: true, error: action.error};
    default:
      return state;
  }
};
