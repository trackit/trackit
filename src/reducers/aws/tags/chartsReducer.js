import Constants from '../../../constants';

export default (state={}, action) => {
  let charts = Object.assign({}, state);
  switch (action.type) {
    case Constants.AWS_TAGS_INSERT_CHARTS:
      return action.charts;
    case Constants.AWS_TAGS_ADD_CHART:
      charts[action.id] = "";
      return charts;
    case Constants.AWS_TAGS_SELECT_KEY:
      charts[action.id] = action.key;
      return charts;
    case Constants.AWS_TAGS_REMOVE_CHART:
      if (charts.hasOwnProperty(action.id))
        delete charts[action.id];
      return charts;
    default:
      return state;
  }
};
