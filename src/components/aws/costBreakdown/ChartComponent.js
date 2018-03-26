import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Spinner from 'react-spinkit';
import BarChart from './BarChartComponent';
import PieChart from './PieChartComponent';
import DifferentiatorChart from './DifferentiatorChartComponent';
import Misc from '../../misc';
import Moment from "moment/moment";

const TimerangeSelector = Misc.TimerangeSelector;
const IntervalNavigator = Misc.IntervalNavigator;
const Selector = Misc.Selector;

/* istanbul ignore next */
const getFilters = (total) => {
  let filters = {
    account: "Account",
    product: "Product",
    region: "Region"
  };
  if (total)
    filters.all = "Total";
  return filters
};

export class Header extends Component {

  constructor(props) {
    super(props);
    this.close = this.close.bind(this);
    this.setDates = this.setDates.bind(this);
    this.setInterval = this.setInterval.bind(this);
    this.setFilter = this.setFilter.bind(this);
  }

  close = (e) => {
    e.preventDefault();
    this.props.close(this.props.id);
  };

  setDates = (start, end) => {
    this.props.setDates(this.props.id, start, end);
  };

  setInterval = (interval) => {
    this.props.setInterval(this.props.id, interval);
  };

  setFilter = (filter) => {
    this.props.setFilter(this.props.id, filter);
  };

  getDateSelector() {
    if (!this.props.setDates)
      return null;
    switch (this.props.type) {
      case "pie":
        return (
          <IntervalNavigator
            startDate={this.props.dates.startDate}
            endDate={this.props.dates.endDate}
            setDatesFunc={this.setDates}
            interval={this.props.interval}
            setIntervalFunc={this.setInterval}
          />
        );
      case "diff":
        return (
          <TimerangeSelector
            startDate={this.props.dates.startDate}
            endDate={this.props.dates.endDate}
            setDatesFunc={this.setDates}
            interval={this.props.interval}
            availableIntervals={["week", "month"]}
            setIntervalFunc={this.setInterval}
            ranges={{
              'Last 2 Weeks': [Moment().subtract(2, 'weeks').startOf('isoWeek'), Moment().subtract(1, 'weeks').endOf('isoWeek')],
              'Last 3 Weeks': [Moment().subtract(3, 'weeks').startOf('isoWeek'), Moment().subtract(1, 'weeks').endOf('isoWeek')],
              'This Month': [Moment().startOf('month'), Moment()],
              'Last Month': [Moment().subtract(1, 'month').startOf('month'), Moment().subtract(1, 'month').endOf('month')],
              'Last 2 Month': [Moment().subtract(2, 'month').startOf('month'), Moment().subtract(1, 'month').endOf('month')],
              'Last 3 Month': [Moment().subtract(3, 'month').startOf('month'), Moment().subtract(1, 'month').endOf('month')],
              'Last 6 Month': [Moment().subtract(3, 'month').startOf('month'), Moment().subtract(1, 'month').endOf('month')],
              'Last 12 Months': [Moment().subtract(1, 'year').startOf('month'), Moment().subtract(1, 'months').endOf('month')],
              'This Year': [Moment().startOf('year'), Moment()],
              'Last Year': [Moment().subtract(1, 'year').startOf('year'), Moment().subtract(1, 'year').endOf('year')]
            }}
          />
        );
      case "bar":
      default:
        return (
          <TimerangeSelector
            startDate={this.props.dates.startDate}
            endDate={this.props.dates.endDate}
            setDatesFunc={this.setDates}
            interval={this.props.interval}
            setIntervalFunc={this.setInterval}
          />
        );
    }
  }

  getIcon() {
    if (!this.props.icon)
      return null;
    switch (this.props.type) {
      case "pie":
        return (
          <div className="cost-breakdown-chart-icon">
            <i className="menu-icon red-color fa fa-pie-chart"/>
            &nbsp;
            Pie Chart
          </div>
        );
      case "diff":
        return (
          <div className="cost-breakdown-chart-icon">
            <i className="menu-icon red-color fa fa-table"/>
            &nbsp;
            Cost Table
          </div>
        );
      case "bar":
        return (
          <div className="cost-breakdown-chart-icon">
            <i className="menu-icon red-color fa fa-pie-chart"/>
            &nbsp;
            Bar Chart
          </div>
        );
      default:
        return null;
    }
  }

  render() {
    const loading = (!this.props.values || !this.props.values.status ? (<Spinner className="spinner clearfix" name='circle'/>) : null);

    const close = (this.props.close ? (
      <button className="btn btn-danger" onClick={this.close}>Remove this chart</button>
    ) : null);

    const error = (this.props.values && this.props.values.status && this.props.values.hasOwnProperty("error") ? (
      <div className="alert alert-warning" role="alert">Data not available ({this.props.values.error.message})</div>
    ) : null);

    const selector = (this.props.type !== "diff" ? (
      <div className="inline-block">
        <Selector
          values={getFilters(!(this.props.type === "pie"))}
          selected={this.props.filter}
          selectValue={this.setFilter}
        />
      </div>
    ) : null);

    const table = (this.props.table && (this.props.type === "pie")? (
      <div className="inline-block table-toggle">
        <button className="btn btn-default" onClick={this.props.toggleTable}>{(this.props.tableStatus ? "Hide" : "Show")} details</button>
      </div>
    ) : null);

    return (
      <div className="clearfix">

        <div className="inline-block pull-left">
          {this.getIcon()}
          {loading}
          {error}
        </div>

        <div className="inline-block pull-right">

          {table}

          {selector}

          <div className="inline-block">
            {this.getDateSelector()}
          </div>

          {close}

        </div>

      </div>
    );
  }

}

Header.propTypes = {
  type: PropTypes.oneOf(["bar", "pie", "diff"]),
  values: PropTypes.object,
  dates: PropTypes.shape({
    startDate: PropTypes.object,
    endDate: PropTypes.object,
  }),
  interval: PropTypes.string.isRequired,
  filter: PropTypes.string.isRequired,
  getCosts: PropTypes.func.isRequired,
  setDates: PropTypes.func,
  setInterval: PropTypes.func.isRequired,
  setFilter: PropTypes.func.isRequired,
  close: PropTypes.func,
  toggleTable: PropTypes.func,
  tableStatus: PropTypes.bool,
  table: PropTypes.bool
};

Header.defaultProps = {
  table: true,
};

class Chart extends Component {

  constructor(props) {
    super(props);
    this.state = {
      table: false
    };
  }

  toggleTable = (e) => {
    e.preventDefault();
    const table = !this.state.table;
    this.setState({table});
  };

  componentWillMount() {
    let filters = [this.props.filter];
    if (this.props.type === "bar")
      filters.push(this.props.interval);
    if (this.props.type === "diff")
      this.props.getCosts(this.props.id, this.props.dates.startDate, this.props.dates.endDate, [this.props.interval], "differentiator");
    else
      this.props.getCosts(this.props.id, this.props.dates.startDate, this.props.dates.endDate, filters, "breakdown");
  }

  componentWillReceiveProps(nextProps) {
    let filters = [nextProps.filter];
    if (nextProps.type === "bar")
      filters.push(nextProps.interval);
    if (this.props.dates !== nextProps.dates ||
      this.props.interval !== nextProps.interval ||
      this.props.filter !== nextProps.filter ||
      this.props.accounts !== nextProps.accounts) {
      if (nextProps.type === "diff")
        nextProps.getCosts(nextProps.id, nextProps.dates.startDate, nextProps.dates.endDate, [nextProps.interval], "differentiator");
      else
        nextProps.getCosts(nextProps.id, nextProps.dates.startDate, nextProps.dates.endDate, filters);
    }
  }

  getChart() {
    if (this.props.values && this.props.values.status && this.props.values.hasOwnProperty("values"))
      switch (this.props.type) {
        case "diff":
          return (<DifferentiatorChart
            values={this.props.values.values}
            interval={this.props.interval}
            legend={this.props.legend}
            height={this.props.height}
            margin={this.props.margin}
          />);
        case "pie":
          return (<PieChart
            values={this.props.values.values}
            interval={this.props.interval}
            filter={this.props.filter}
            legend={this.props.legend}
            height={this.props.height}
            margin={this.props.margin}
            table={this.state.table}
          />);
        case "bar":
        default:
          return (<BarChart
            values={this.props.values.values}
            interval={this.props.interval}
            filter={this.props.filter}
            legend={this.props.legend}
            height={this.props.height}
            margin={this.props.margin}
          />);
      }
      return (<div className="no-chart" style={{height: this.props.height}}/>);
  }

  render() {
    const chart = this.getChart();

    return (
      <div className="clearfix">
        <Header {...this.props} toggleTable={this.toggleTable} tableStatus={this.state.table}/>
        {chart}
      </div>
    );
  }

}

Chart.propTypes = {
  id: PropTypes.string.isRequired,
  type: PropTypes.oneOf(["bar", "pie", "diff"]),
  values: PropTypes.object,
  dates: PropTypes.shape({
    startDate: PropTypes.object,
    endDate: PropTypes.object,
  }),
  accounts: PropTypes.arrayOf(PropTypes.object),
  interval: PropTypes.string.isRequired,
  filter: PropTypes.string.isRequired,
  getCosts: PropTypes.func.isRequired,
  setDates: PropTypes.func,
  setInterval: PropTypes.func.isRequired,
  setFilter: PropTypes.func.isRequired,
  close: PropTypes.func,
  legend: PropTypes.bool,
  height: PropTypes.number,
  margin: PropTypes.bool,
  table: PropTypes.bool,
  icon: PropTypes.bool
};

Chart.defaultProps = {
  legend: true,
  height: 400,
  margin: true,
  table: true,
  icon: true
};

export default Chart;
