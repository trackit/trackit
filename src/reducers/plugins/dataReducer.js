import Constants from '../../constants';

export default (state={}, action) => {
  switch (action.type) {
    case Constants.GET_PLUGINS_DATA:
      return { status: false };
    case Constants.GET_PLUGINS_DATA_SUCCESS:
      return { status: true, values: action.data };
    case Constants.GET_PLUGINS_DATA_ERROR:
      return { status: true, error: action.error };
    default:
      return state;
  }
};
