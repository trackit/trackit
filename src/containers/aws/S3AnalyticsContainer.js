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
  startDate: moment().startOf('month'),
  endDate: moment()
};

// S3AnalyticsContainer Component
export class S3AnalyticsContainer extends Component {

  componentDidMount() {
    this.props.setDates(defaultDates.startDate, defaultDates.endDate);
    if (this.props.dates)
      this.props.getData(this.props.dates.startDate, this.props.dates.endDate);
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.dates !== nextProps.dates)
      nextProps.getData(nextProps.dates.startDate, nextProps.dates.endDate);
  }

  render() {
    return (
      <Panel>

        <div className="clearfix">
          <h3 className="white-box-title no-padding inline-block">
            <img className="white-box-title-icon" src={s3square} alt="AWS square logo"/>
            AWS S3 Analytics
          </h3>
          <div className="inline-block pull-right">
            <TimerangeSelector
              startDate={this.props.dates.startDate}
              endDate={this.props.dates.endDate}
              setDatesFunc={this.props.setDates}
            />
          </div>
        </div>

        {this.props.values && <S3Analytics.Infos data={this.props.values}/>}

        {this.props.values && <S3Analytics.BarChart elementId="s3BarChart" data={this.props.values}/>}

        <div className="no-padding">
          {this.props.values && <S3Analytics.Table data={this.props.values}/>}
        </div>

      </Panel>
    );
  }

}

S3AnalyticsContainer.propTypes = {
  values: PropTypes.object,
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
