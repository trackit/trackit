import Constants from '../../constants';

export default (state={}, action) => {
  let items = Object.assign({}, state);
  switch (action.type) {
    case Constants.DASHBOARD_INSERT_ITEMS:
    case Constants.DASHBOARD_UPDATE_ITEMS:
      return action.items;
    case Constants.DASHBOARD_ADD_ITEM:
      items[action.id] = action.props;
      return items;
    case Constants.DASHBOARD_REMOVE_ITEM:
      if (items.hasOwnProperty(action.id))
        delete items[action.id];
      return items;
    default:
      return state;
  }
};
