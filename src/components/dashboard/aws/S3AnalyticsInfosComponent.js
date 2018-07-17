import React, { Component } from 'react';
import PropTypes from 'prop-types';
import AWS from '../../aws';

const Infos = AWS.S3Analytics.Infos;

class S3AnalyticsInfosComponent extends Component {

  componentWillMount() {
    if (this.props.dates)
      this.props.getValues(this.props.id, "s3", this.props.dates.startDate, this.props.dates.endDate, null);
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.dates !== nextProps.dates ||
      this.props.accounts !== nextProps.accounts)
      nextProps.getValues(nextProps.id, "s3", nextProps.dates.startDate, nextProps.dates.endDate, null);
  }

  render() {
    return (
      <Infos data={this.props.values} offset={false}/>
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
};

export default S3AnalyticsInfosComponent;
