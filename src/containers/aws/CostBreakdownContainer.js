import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import UUID from 'uuid/v4';
import Components from '../../components';
import Actions from '../../actions';
import s3square from '../../assets/s3-square.png';

const Panel = Components.Misc.Panel;
const Chart = Components.AWS.CostBreakdown.Chart;
const Infos = Components.AWS.CostBreakdown.Infos;

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

export class CostBreakdownContainer extends Component {

  constructor(props) {
    super(props);
    if (!this.props.charts || !Object.keys(this.props.charts).length)
      this.props.initCharts();
    this.addBarChart = this.addBarChart.bind(this);
    this.addPieChart = this.addPieChart.bind(this);
    this.resetCharts = this.resetCharts.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    if (!Object.keys(nextProps.charts).length)
      nextProps.initCharts();
    clearTooltips();
  }

  addSummary = (e) => {
    e.preventDefault();
    this.props.addChart("summary");
  };

  addBarChart = (e) => {
    e.preventDefault();
    this.props.addChart("bar");
  };

  addPieChart = (e) => {
    e.preventDefault();
    this.props.addChart("pie");
  };

  addDiffChart = (e) => {
    e.preventDefault();
    this.props.addChart("diff");
  };

  resetCharts = (e) => {
    e.preventDefault();
    Object.keys(this.props.charts).forEach((id) => {this.props.removeChart(id)});
  };

  getChart(id, chartType, index) {
    if (this.props.costsValues &&
      this.props.costsDates && this.props.costsDates.hasOwnProperty(id) &&
      this.props.costsInterval && this.props.costsInterval.hasOwnProperty(id) &&
      this.props.costsFilter && this.props.costsFilter.hasOwnProperty(id)
    ) {
      if (chartType === "summary")
        return (
          <Infos
            key={index}
            id={id}
            accounts={this.props.accounts}
            values={this.props.costsValues[id]}
            dates={this.props.costsDates[id]}
            interval={this.props.costsInterval[id]}
            getCosts={this.props.getCosts}
            setDates={this.props.setCostsDates}
            setInterval={this.props.setCostsInterval}
            close={Object.keys(this.props.charts).length > minimalCount ? this.props.removeChart : null}
          />
        );
      return (
        <Chart
          key={index}
          id={id}
          type={chartType}
          accounts={this.props.accounts}
          values={this.props.costsValues[id]}
          dates={this.props.costsDates[id]}
          interval={this.props.costsInterval[id]}
          filter={this.props.costsFilter[id]}
          getCosts={this.props.getCosts}
          setDates={this.props.setCostsDates}
          setInterval={this.props.setCostsInterval}
          setFilter={this.props.setCostsFilter}
          close={Object.keys(this.props.charts).length > minimalCount ? this.props.removeChart : null}
        />
      );
    }
    return null;
  }

  render() {
    const header = (
      <div className="clearfix">
        <h3 className="white-box-title no-padding inline-block">
          <img className="white-box-title-icon" src={s3square} alt="AWS square logo"/>
          Cost Breakdown
        </h3>
        <div className="inline-block pull-right">
          <button className="btn btn-default inline-block" onClick={this.addSummary}>
            <i className="menu-icon fa fa-list"/>
            &nbsp;
            Add a summary
          </button>
          &nbsp;
          <button className="btn btn-default inline-block" onClick={this.addBarChart}>
            <i className="menu-icon fa fa-bar-chart"/>
            &nbsp;
            Add a bar chart
          </button>
          &nbsp;
          <button className="btn btn-default inline-block" onClick={this.addPieChart}>
            <i className="menu-icon fa fa-pie-chart"/>
            &nbsp;
            Add a pie chart
          </button>
          &nbsp;
          <button className="btn btn-default inline-block" onClick={this.addDiffChart}>
            <i className="menu-icon fa fa-table"/>
            &nbsp;
            Add a cost table
          </button>
          &nbsp;
          <button className="btn btn-danger inline-block" onClick={this.resetCharts}>Reset charts</button>
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

CostBreakdownContainer.propTypes = {
  costsValues: PropTypes.object,
  costsDates: PropTypes.object,
  charts: PropTypes.object,
  accounts: PropTypes.arrayOf(PropTypes.object),
  costsInterval: PropTypes.object.isRequired,
  costsFilter: PropTypes.object.isRequired,
  initCharts: PropTypes.func.isRequired,
  addChart: PropTypes.func.isRequired,
  removeChart: PropTypes.func.isRequired,
  getCosts: PropTypes.func.isRequired,
  setCostsDates: PropTypes.func.isRequired,
  setCostsInterval: PropTypes.func.isRequired,
  setCostsFilter: PropTypes.func.isRequired,
  resetCostsDates: PropTypes.func.isRequired,
  resetCostsInterval: PropTypes.func.isRequired,
  resetCostsFilter: PropTypes.func.isRequired,
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  charts: aws.costs.charts,
  costsValues: aws.costs.values,
  costsDates: aws.costs.dates,
  costsInterval: aws.costs.interval,
  costsFilter: aws.costs.filter,
  accounts: aws.accounts.selection
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  initCharts: () => {
    dispatch(Actions.AWS.Costs.initCharts());
  },
  addChart: (type) => {
    dispatch(Actions.AWS.Costs.addChart(UUID(), type));
  },
  removeChart: (id) => {
    dispatch(Actions.AWS.Costs.removeChart(id));
  },
  getCosts: (id, begin, end, filters, chartType) => {
    dispatch(Actions.AWS.Costs.getCosts(id, begin, end, filters, chartType));
  },
  setCostsDates: (id, startDate, endDate) => {
    dispatch(Actions.AWS.Costs.setCostsDates(id, startDate, endDate))
  },
  resetCostsDates: () => {
    dispatch(Actions.AWS.Costs.resetCostsDates())
  },
  setCostsInterval: (id, interval) => {
    dispatch(Actions.AWS.Costs.setCostsInterval(id, interval));
  },
  resetCostsInterval: () => {
    dispatch(Actions.AWS.Costs.resetCostsInterval());
  },
  setCostsFilter: (id, filter) => {
    dispatch(Actions.AWS.Costs.setCostsFilter(id, filter));
  },
  resetCostsFilter: () => {
    dispatch(Actions.AWS.Costs.resetCostsFilter());
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(CostBreakdownContainer);
