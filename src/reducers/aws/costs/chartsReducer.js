import Constants from '../../../constants';

export default (state=[], action) => {
  let charts = Object.assign([], state);
  switch (action.type) {
    case Constants.AWS_INSERT_CHARTS:
      return action.charts;
    case Constants.AWS_ADD_CHART:
      charts.push(action.id);
      return charts;
    case Constants.AWS_REMOVE_CHART:
      const index = charts.indexOf(action.id);
      if (index > -1)
        charts.splice(index, 1);
      return charts;
    default:
      return state;
  }
};
