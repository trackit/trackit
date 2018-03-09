import Constants from '../../constants';

export default (state={}, action) => {
  let values = Object.assign({}, state);
  switch (action.type) {
    case Constants.DASHBOARD_GET_VALUES_ERROR:
      values[action.id] = { status: true, error: action.error };
      return values;
    case Constants.DASHBOARD_GET_VALUES:
      values[action.id] = { status: false };
      return values;
    case Constants.DASHBOARD_GET_VALUES_SUCCESS:
      values[action.id] = { status: true, values: action.data };
      return values;
    case Constants.DASHBOARD_REMOVE_ITEM:
      if (values.hasOwnProperty(action.id))
        delete values[action.id];
      return values;
    default:
      return state;
  }
};
