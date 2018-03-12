import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Spinner from 'react-spinkit';
import AWS from '../../aws';
import Misc from '../../misc';

const TimerangeSelector = Misc.TimerangeSelector;
const Selector = Misc.Selector;
const Charts = AWS.S3Analytics;

/* istanbul ignore next */
const filters = {
  storage: "Storage",
  bandwidth: "Bandwidth",
  requests: "Requests"
};

class S3AnalyticsChartsComponent extends Component {

  constructor(props) {
    super(props);
    this.setDates = this.setDates.bind(this);
    this.setFilter = this.setFilter.bind(this);
  }

  componentWillMount() {
    if (this.props.dates)
      this.props.getValues(this.props.id, "s3", this.props.dates.startDate, this.props.dates.endDate, null);
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.dates !== nextProps.dates ||
      this.props.filter !== nextProps.filter ||
      this.props.accounts !== nextProps.accounts)
      nextProps.getValues(nextProps.id, "s3", nextProps.dates.startDate, nextProps.dates.endDate, null);
  }

  setDates = (startDate, endDate) => {
    this.props.setDates(this.props.id, startDate, endDate);
  };

  setFilter = (filter) => {
    this.props.setFilter(this.props.id, filter);
  };

  getChart() {
    if (!this.props.values || !this.props.values.status)
      return (<Spinner className="spinner clearfix" name='circle'/>);
    switch (this.props.filter) {
      case "storage":
        return (<Charts.StorageCostChart data={this.props.values}/>);
      case "bandwidth":
        return (<Charts.BandwidthCostChart data={this.props.values}/>);
      case "requests":
        return (<Charts.RequestsCostChart data={this.props.values}/>);
      default:
        return null
    }
  }

  render() {
    const timerange = (this.props.dates ? (
      <TimerangeSelector
        startDate={this.props.dates.startDate}
        endDate={this.props.dates.endDate}
        setDatesFunc={this.setDates}
      />
    ) : null);

    return (
      <div>
        <div className="clearfix">
          <div className="inline-block pull-right">
            <Selector
              values={filters}
              selected={this.props.filter}
              selectValue={this.setFilter}
            />
            {timerange}
          </div>
        </div>
        {this.getChart()}
      </div>
    );
  }

}

S3AnalyticsChartsComponent.propTypes = {
  id: PropTypes.string.isRequired,
  accounts: PropTypes.arrayOf(PropTypes.object),
  values: PropTypes.object,
  getValues: PropTypes.func.isRequired,
  dates: PropTypes.shape({
    startDate: PropTypes.object,
    endDate: PropTypes.object,
  }),
  setDates: PropTypes.func.isRequired,
  filter: PropTypes.string.isRequired,
  setFilter: PropTypes.func.isRequired
};

export default S3AnalyticsChartsComponent;
