import React, {Component} from 'react';
import {connect} from 'react-redux';
// import PropTypes from 'prop-types';

import Actions from '../../actions';

import Components from '../../components';

import s3square from '../../assets/s3-square.png';
import PropTypes from "prop-types";

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
          <h3 className="white-box-title no-padding">
            <img className="white-box-title-icon" src={s3square} alt="AWS square logo"/>
            AWS S3 Analytics
          </h3>
        </div>

        <div className="row">
          <div className="col-md-12">
            <div className="white-box">
              {this.props.s3Data && <S3Analytics.Infos data={this.props.s3Data}/>}
            </div>
          </div>
          <div className="col-md-12">
            <div className="white-box">
              {this.props.s3Data.length && <S3Analytics.BarChart elementId="s3BarChart" data={this.props.s3Data}/>}
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
  getS3Data: PropTypes.func.isRequired
};

const mapStateToProps = ({aws}) => ({
  s3Data: aws.s3
});

const mapDispatchToProps = (dispatch) => ({
  getS3Data: () => {
    dispatch(Actions.AWS.S3.getS3Data())
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(S3AnalyticsContainer);
