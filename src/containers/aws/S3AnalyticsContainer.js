import React, {Component} from 'react';
import {connect} from 'react-redux';
import PropTypes from 'prop-types';
import moment from "moment/moment";

import Actions from '../../actions';

import Components from '../../components';

import s3square from '../../assets/s3-square.png';

const TimerangeSelector = Components.Misc.TimerangeSelector;
const S3Analytics = Components.AWS.S3Analytics;
const Panel = Components.Misc.Panel;

const defaultDates = {
  startDate: moment().subtract(1, 'months').startOf('month'),
  endDate: moment().subtract(1, 'months').endOf('month')
};

// S3AnalyticsContainer Component
export class S3AnalyticsContainer extends Component {

  componentDidMount() {
    if (this.props.dates)
      this.props.getData(this.props.dates.startDate, this.props.dates.endDate);
    else
      this.props.setDates(defaultDates.startDate, defaultDates.endDate);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.dates && (this.props.dates !== nextProps.dates || this.props.accounts !== nextProps.accounts))
      nextProps.getData(nextProps.dates.startDate, nextProps.dates.endDate);
  }

  render() {
    const timerange = (this.props.dates ?  (
      <TimerangeSelector
        startDate={this.props.dates.startDate}
        endDate={this.props.dates.endDate}
        setDatesFunc={this.props.setDates}
      />
    ) : null);

    return (
      <Panel>

        <div className="clearfix">
          <h3 className="white-box-title no-padding inline-block">
            <img className="white-box-title-icon" src={s3square} alt="AWS square logo"/>
            AWS S3 Analytics
          </h3>
          <div className="inline-block pull-right">
            {timerange}
          </div>
        </div>

        <S3Analytics.Infos data={this.props.values}/>

        <div>
          <div className="row">
            <div className="col-md-6">
              <S3Analytics.BandwidthCostChart data={this.props.values}/>
            </div>
            <div className="col-md-6">
              <S3Analytics.StorageCostChart data={this.props.values}/>
            </div>
          </div>
        </div>

        <div className="no-padding">
          <S3Analytics.Table data={this.props.values}/>
        </div>

      </Panel>
    );
  }

}

S3AnalyticsContainer.propTypes = {
  values: PropTypes.object,
  accounts: PropTypes.arrayOf(PropTypes.object),
  dates: PropTypes.shape({
    startDate: PropTypes.object.isRequired,
    endDate: PropTypes.object.isRequired,
  }),
  getData: PropTypes.func.isRequired,
  setDates: PropTypes.func.isRequired
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  values: aws.s3.values,
  dates: aws.s3.dates,
  accounts: aws.accounts.selection
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getData: (begin, end) => {
    dispatch(Actions.AWS.S3.getData(begin, end))
  },
  setDates: (startDate, endDate) => {
    dispatch(Actions.AWS.S3.setDates(startDate, endDate))
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(S3AnalyticsContainer);
