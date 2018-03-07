import React, { Component } from 'react';
import PropTypes from 'prop-types';
import AWS from '../../aws';
import Misc from '../../misc';

const TimerangeSelector = Misc.TimerangeSelector;
const Infos = AWS.S3Analytics.Infos;

class S3AnalyticsInfosComponent extends Component {

  constructor(props) {
    super(props);
    this.setDates = this.setDates.bind(this);
  }

  componentWillMount() {
    if (this.props.dates)
      this.props.getValues(this.props.id, "s3", this.props.dates.startDate, this.props.dates.endDate, null);
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.dates !== nextProps.dates ||
      this.props.accounts !== nextProps.accounts)
      nextProps.getValues(nextProps.id, "s3", nextProps.dates.startDate, nextProps.dates.endDate, null);
  }

  setDates = (startDate, endDate) => {
    this.props.setDates(this.props.id, startDate, endDate);
  };

  render() {
    const timerange = (this.props.dates ?  (
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
            {timerange}
          </div>
        </div>
        <Infos data={this.props.values}/>
      </div>
    );
  }

}

S3AnalyticsInfosComponent.propTypes = {
  id: PropTypes.string.isRequired,
  accounts: PropTypes.arrayOf(PropTypes.object),
  values: PropTypes.object,
  getValues: PropTypes.func.isRequired,
  dates: PropTypes.shape({
    startDate: PropTypes.object,
    endDate: PropTypes.object,
  }),
  setDates: PropTypes.func.isRequired,
};

export default S3AnalyticsInfosComponent;
