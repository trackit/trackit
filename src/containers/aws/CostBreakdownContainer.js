import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import UUID from 'uuid/v4';
import Components from '../../components';
import Actions from '../../actions';
import s3square from '../../assets/s3-square.png';

const TimerangeSelector = Components.Misc.TimerangeSelector;
const Selector = Components.Misc.Selector;
const Panel = Components.Misc.Panel;
const CostBreakdownChart = Components.AWS.CostBreakdown.Chart;

const filters = {
  all: "Total",
  account: "Account",
  product: "Product",
  region: "Region"
};

export class Chart extends Component {

  constructor(props) {
    super(props);
    this.setDates = this.setDates.bind(this);
    this.setInterval = this.setInterval.bind(this);
    this.setFilter = this.setFilter.bind(this);
    this.close = this.close.bind(this);
  }

  componentWillMount() {
    this.props.getCosts(this.props.id, this.props.dates.startDate, this.props.dates.endDate, [this.props.filter, this.props.interval]);
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.dates !== nextProps.dates ||
      this.props.interval !== nextProps.interval ||
      this.props.filter !== nextProps.filter ||
      this.props.accounts !== nextProps.accounts)
      nextProps.getCosts(nextProps.id, nextProps.dates.startDate, nextProps.dates.endDate, [nextProps.filter, nextProps.interval]);
  }

  setDates = (start, end) => {
    this.props.setDates(this.props.id, start, end);
  };


  setInterval = (interval) => {
    this.props.setInterval(this.props.id, interval);
  };


  setFilter = (filter) => {
    this.props.setFilter(this.props.id, filter);
  };

  close = (e) => {
    e.preventDefault();
    this.props.close(this.props.id);
  };

  render() {
    const close = (this.props.close ? (
      <button className="btn btn-danger" onClick={this.close}>Remove this chart</button>
    ) : null);
    return (
      <div className="clearfix">
        <div className="inline-block pull-right">
          <div className="inline-block">
            <Selector
              values={filters}
              selected={this.props.filter}
              selectValue={this.setFilter}
            />
          </div>
          <div className="inline-block">
            <TimerangeSelector
              startDate={this.props.dates.startDate}
              endDate={this.props.dates.endDate}
              setDatesFunc={this.setDates}
              interval={this.props.interval}
              setIntervalFunc={this.setInterval}
            />
          </div>
          {close}
        </div>
        <CostBreakdownChart values={this.props.values} interval={this.props.interval} filter={this.props.filter}/>
      </div>
    );
  }

}

Chart.propTypes = {
  id: PropTypes.string.isRequired,
  values: PropTypes.object,
  dates: PropTypes.shape({
    startDate: PropTypes.object,
    endDate: PropTypes.object,
  }),
  accounts: PropTypes.arrayOf(PropTypes.string),
  interval: PropTypes.string.isRequired,
  filter: PropTypes.string.isRequired,
  getCosts: PropTypes.func.isRequired,
  setDates: PropTypes.func.isRequired,
  setInterval: PropTypes.func.isRequired,
  setFilter: PropTypes.func.isRequired,
  close: PropTypes.func
};

export class CostBreakdownContainer extends Component {

  constructor(props) {
    super(props);
    if (!this.props.charts || !this.props.charts.length)
      this.props.addChart();
    this.addChart = this.addChart.bind(this);
    this.resetCharts = this.resetCharts.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    if (!nextProps.charts.length)
      nextProps.addChart();
  }

  addChart = (e) => {
    e.preventDefault();
    this.props.addChart();
  };

  resetCharts = (e) => {
    e.preventDefault();
    this.props.charts.forEach((id) => {this.props.removeChart(id)});
  };

  getChart(id, index) {
    if (this.props.costsValues &&
      this.props.costsDates && this.props.costsDates.hasOwnProperty(id) &&
      this.props.costsInterval && this.props.costsInterval.hasOwnProperty(id) &&
      this.props.costsFilter && this.props.costsFilter.hasOwnProperty(id)
    )
      return (
        <Chart
          key={index}
          id={id}
          accounts={this.props.accounts}
          values={this.props.costsValues[id]}
          dates={this.props.costsDates[id]}
          interval={this.props.costsInterval[id]}
          filter={this.props.costsFilter[id]}
          getCosts={this.props.getCosts}
          setDates={this.props.setCostsDates}
          setInterval={this.props.setCostsInterval}
          setFilter={this.props.setCostsFilter}
          close={this.props.charts.length > 1 ? this.props.removeChart : null}
        />
      );
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
          <button className="btn btn-default inline-block" onClick={this.addChart}>Add a chart</button>
          &nbsp;
          <button className="btn btn-danger inline-block" onClick={this.resetCharts}>Reset charts</button>
        </div>
      </div>
    );
    const charts = this.props.charts.map((id, index) => (this.getChart(id, index)));
    const children = [header, ...charts];
    return(
      <Panel children={children}/>
    );
  }
}

CostBreakdownContainer.propTypes = {
  costsValues: PropTypes.object,
  costsDates: PropTypes.object,
  charts: PropTypes.arrayOf(PropTypes.string),
  accounts: PropTypes.arrayOf(PropTypes.string),
  costsInterval: PropTypes.object.isRequired,
  costsFilter: PropTypes.object.isRequired,
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
  addChart: () => {
    dispatch(Actions.AWS.Costs.addChart(UUID()));
  },
  removeChart: (id) => {
    dispatch(Actions.AWS.Costs.removeChart(id));
  },
  getCosts: (id, begin, end, filters) => {
    dispatch(Actions.AWS.Costs.getCosts(id, begin, end, filters));
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
