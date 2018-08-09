import React, {Component} from 'react';

import PropTypes from "prop-types";
import {connect} from "react-redux";
import UUID from "uuid/v4";
import Components from '../../components';
import Actions from '../../actions';

const Panel = Components.Misc.Panel;
const Chart = Components.AWS.Tags.Chart;

// This function will hide NVD3 tooltips to avoid ghost tooltips to stay on screen when chart they are linked to is updated or deleted
// Similar issue : https://github.com/novus/nvd3/issues/1262
/* istanbul ignore next */
const clearTooltips = () => {
  const tooltips = document.getElementsByClassName("nvtooltip xy-tooltip");
  for (let i = 0; i < tooltips.length; i++) {
    tooltips[i].style.opacity = 0;
  }
};

const minimalCount = 2;

export class TagsContainer extends Component {

  constructor(props) {
    super(props);
    if (!this.props.charts || !Object.keys(this.props.charts).length)
      this.props.initCharts();
    this.addChart = this.addChart.bind(this);
    this.resetCharts = this.resetCharts.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    if (!Object.keys(nextProps.charts).length)
      nextProps.initCharts();
    clearTooltips();
  }

  componentWillUnmount() {
    this.props.clearKeys();
  }

  addChart = (e) => {
    e.preventDefault();
    this.props.addChart();
  };

  resetCharts = (e) => {
    e.preventDefault();
    Object.keys(this.props.charts).forEach((id) => {this.props.removeChart(id)});
  };

  getChart(id, tag, index) {
    if (this.props.values
      && this.props.dates && this.props.dates.hasOwnProperty(id)
      && this.props.interval && this.props.interval.hasOwnProperty(id)) {
      return (
        <Chart
          key={index}
          id={id}
          accounts={this.props.accounts}
          keys={this.props.keys[id]}
          getKeys={this.props.getKeys}
          tag={tag}
          selectKey={this.props.selectKey}
          values={this.props.values[id]}
          getValues={this.props.getValues}
          dates={this.props.dates[id]}
          setDates={this.props.setDates}
          interval={this.props.interval[id]}
          setInterval={this.props.setInterval}
          close={Object.keys(this.props.charts).length > minimalCount ? this.props.removeChart : null}
        />
      );
    }
    return null;
  }

  render() {
    const header = (
      <div className="clearfix">
        <h3 className="white-box-title no-padding p-l-15 inline-block">
          <i className="fa fa-tags"/>
          &nbsp;
          Tags
        </h3>
        <div className="inline-block pull-right">
          <button className="btn btn-default inline-block" onClick={this.addChart}>
            <i className="menu-icon fa fa-pie-chart"/>
            &nbsp;
            Add a panel
          </button>
          &nbsp;
          <button className="btn btn-danger inline-block" onClick={this.resetCharts}>
            <i className="fa fa-eraser"/>
            &nbsp;
            Reset panels
          </button>
        </div>
      </div>
    );
    const charts = Object.keys(this.props.charts).map((id, index) => (this.getChart(id, this.props.charts[id], index)));
    const children = [header, ...charts];
    return(
      <Panel children={children}/>
    );
  }
}

TagsContainer.propTypes = {
  accounts: PropTypes.arrayOf(PropTypes.object),
  charts: PropTypes.object.isRequired,
  keys: PropTypes.object.isRequired,
  values: PropTypes.object.isRequired,
  dates: PropTypes.object.isRequired,
  interval: PropTypes.object.isRequired,
  initCharts: PropTypes.func.isRequired,
  addChart: PropTypes.func.isRequired,
  removeChart: PropTypes.func.isRequired,
  getValues: PropTypes.func.isRequired,
  setDates: PropTypes.func.isRequired,
  resetDates: PropTypes.func.isRequired,
  clearDates: PropTypes.func.isRequired,
  getKeys: PropTypes.func.isRequired,
  clearKeys: PropTypes.func.isRequired,
  selectKey: PropTypes.func.isRequired,
  setInterval: PropTypes.func.isRequired,
  resetInterval: PropTypes.func.isRequired,
  clearInterval: PropTypes.func.isRequired
};

/* istanbul ignore next */
const tagsStateToProps = ({aws}) => ({
  charts: aws.tags.charts,
  keys: aws.tags.keys,
  dates: aws.tags.dates,
  interval: aws.tags.interval,
  values: aws.tags.values,
  accounts: aws.accounts.selection
});

/* istanbul ignore next */
const tagsDispatchToProps = (dispatch) => ({
  initCharts: () => {
    dispatch(Actions.AWS.Tags.initCharts());
  },
  addChart: () => {
    dispatch(Actions.AWS.Tags.addChart(UUID()));
  },
  removeChart: (id) => {
    dispatch(Actions.AWS.Tags.removeChart(id));
  },
  getValues: (id, begin, end, key) => {
    dispatch(Actions.AWS.Tags.getValues(id, begin, end, key));
  },
  setDates: (id, startDate, endDate) => {
    dispatch(Actions.AWS.Tags.setDates(id, startDate, endDate))
  },
  resetDates: () => {
    dispatch(Actions.AWS.Tags.resetCostsDates())
  },
  clearDates: () => {
    dispatch(Actions.AWS.Tags.clearDates())
  },
  getKeys: (id, start, end) => {
    dispatch(Actions.AWS.Tags.getKeys(id, start, end));
  },
  clearKeys: () => {
    dispatch(Actions.AWS.Tags.clearKeys());
  },
  selectKey: (id, tag) => {
    dispatch(Actions.AWS.Tags.selectKey(id, tag));
  },
  setInterval: (id, interval) => {
    dispatch(Actions.AWS.Tags.setInterval(id, interval))
  },
  resetInterval: () => {
    dispatch(Actions.AWS.Tags.resetInterval())
  },
  clearInterval: () => {
    dispatch(Actions.AWS.Tags.clearInterval())
  },
});

export default connect(tagsStateToProps, tagsDispatchToProps)(TagsContainer);
