import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Spinner from 'react-spinkit';
import BarChart from './BarChartComponent';
import PieChart from './PieChartComponent';
import Misc from '../../misc';

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

  render() {
    const loading = (!this.props.values || !this.props.values.status ? (<Spinner className="spinner clearfix" name='circle'/>) : null);

    const close = (this.props.close ? (
      <button className="btn btn-danger" onClick={this.close}>Remove this chart</button>
    ) : null);

    const error = (this.props.values && this.props.values.status && this.props.values.hasOwnProperty("error") ? (
      <div className="alert alert-warning" role="alert">Data not available ({this.props.values.error.message})</div>
    ) : null);

    return (
      <div>

        <div className="inline-block pull-left">
          {loading}
          {error}
        </div>

        <div className="inline-block pull-right">

          <div className="inline-block">
            <Selector
              values={getFilters(!(this.props.type === "pie"))}
              selected={this.props.filter}
              selectValue={this.setFilter}
            />
          </div>

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
  type: PropTypes.oneOf(["bar", "pie"]),
  values: PropTypes.object,
  dates: PropTypes.shape({
    startDate: PropTypes.object,
    endDate: PropTypes.object,
  }),
  interval: PropTypes.string.isRequired,
  filter: PropTypes.string.isRequired,
  getCosts: PropTypes.func.isRequired,
  setDates: PropTypes.func.isRequired,
  setInterval: PropTypes.func.isRequired,
  setFilter: PropTypes.func.isRequired,
  close: PropTypes.func
};

class Chart extends Component {

  componentWillMount() {
    let filters = [this.props.filter];
    if (this.props.type === "bar")
      filters.push(this.props.interval);
    this.props.getCosts(this.props.id, this.props.dates.startDate, this.props.dates.endDate, filters);
  }

  componentWillReceiveProps(nextProps) {
    let filters = [nextProps.filter];
    if (nextProps.type === "bar")
      filters.push(nextProps.interval);
    if (this.props.dates !== nextProps.dates ||
      this.props.interval !== nextProps.interval ||
      this.props.filter !== nextProps.filter ||
      this.props.accounts !== nextProps.accounts)
        nextProps.getCosts(nextProps.id, nextProps.dates.startDate, nextProps.dates.endDate, filters);
  }

  getChart() {
    if (this.props.values && this.props.values.status && this.props.values.hasOwnProperty("values"))
      switch (this.props.type) {
        case "pie":
          return (<PieChart
            values={this.props.values.values}
            interval={this.props.interval}
            filter={this.props.filter}
            legend={this.props.legend}
            height={this.props.height}
            margin={this.props.margin}
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
  }

  render() {
    const chart = this.getChart();

    return (
      <div className="clearfix">
        <Header {...this.props} />
        {chart}
      </div>
    );
  }

}

Chart.propTypes = {
  id: PropTypes.string.isRequired,
  type: PropTypes.oneOf(["bar", "pie"]),
  values: PropTypes.object,
  dates: PropTypes.shape({
    startDate: PropTypes.object,
    endDate: PropTypes.object,
  }),
  accounts: PropTypes.arrayOf(PropTypes.object),
  interval: PropTypes.string.isRequired,
  filter: PropTypes.string.isRequired,
  getCosts: PropTypes.func.isRequired,
  setDates: PropTypes.func.isRequired,
  setInterval: PropTypes.func.isRequired,
  setFilter: PropTypes.func.isRequired,
  close: PropTypes.func,
  legend: PropTypes.bool,
  height: PropTypes.number,
  margin: PropTypes.bool
};

Chart.defaultProps = {
  legend: true,
  height: 400,
  margin: true
};

export default Chart;
