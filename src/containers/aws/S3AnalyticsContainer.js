import React, {Component} from 'react';
import {connect} from 'react-redux';
import PropTypes from 'prop-types';

import Actions from '../../actions';

import Components from '../../components';

import s3square from '../../assets/s3-square.png';

const TimerangeSelector = Components.Misc.TimerangeSelector;
const S3Analytics = Components.AWS.S3Analytics;

// S3AnalyticsContainer Component
export class S3AnalyticsContainer extends Component {

  componentDidMount() {
    this.props.getS3Data();
  }

  render() {
    return (
      <div className="container-fluid">

        <div className="white-box">
          <h3 className="white-box-title no-padding inline-block">
            <img className="white-box-title-icon" src={s3square} alt="AWS square logo"/>
            AWS S3 Analytics
          </h3>
          <div className="inline-block pull-right">
            <TimerangeSelector
              startDate={this.props.s3View.startDate}
              endDate={this.props.s3View.endDate}
              setDatesFunc={this.props.setS3ViewDates}
            />
          </div>
        </div>

        <div className="row">
          <div className="col-md-12">
            <div className="white-box">
              {this.props.s3Data && <S3Analytics.Infos data={this.props.s3Data}/>}
            </div>
          </div>
          <div className="col-md-12">
            <div className="white-box">
              {this.props.s3Data && <S3Analytics.BarChart elementId="s3BarChart" data={this.props.s3Data}/>}
            </div>
          </div>

        </div>

        <div className="white-box no-padding">
          {this.props.s3Data && <S3Analytics.Table data={this.props.s3Data}/>}
        </div>
      </div>
    );
  }

}

S3AnalyticsContainer.propTypes = {
  s3Data: PropTypes.arrayOf(
    PropTypes.shape({
      _id: PropTypes.string.isRequired,
      size: PropTypes.number.isRequired,
      storage_cost: PropTypes.number.isRequired,
      bw_cost: PropTypes.number.isRequired,
      total_cost: PropTypes.number.isRequired,
      transfer_in: PropTypes.number.isRequired,
      transfer_out: PropTypes.number.isRequired,
    })
  ),
  s3View: PropTypes.shape({
    startDate: PropTypes.object.isRequired,
    endDate: PropTypes.object.isRequired,
  }),
  getS3Data: PropTypes.func.isRequired,
  setS3ViewDates: PropTypes.func.isRequired
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  s3Data: aws.s3.data,
  s3View: aws.s3.view,
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getS3Data: () => {
    dispatch(Actions.AWS.S3.getS3Data())
  },
  setS3ViewDates: (startDate, endDate) => {
    dispatch(Actions.AWS.S3.setS3ViewDates(startDate, endDate))
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(S3AnalyticsContainer);
