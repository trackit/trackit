import React, {Component} from 'react';
import {connect} from 'react-redux';
// import PropTypes from 'prop-types';

import Actions from '../actions';

import Components from '../components';

import s3square from '../assets/s3-square.png';

// S3AnalyticsContainer Component
class S3AnalyticsContainer extends Component {

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
              {this.props.s3Data && <Components.S3Analytics.S3AnalyticsInfos data={this.props.s3Data}/>}
            </div>
          </div>
          <div className="col-md-12">
            <div className="white-box">
              {this.props.s3Data && <Components.S3Analytics.S3AnalyticsBarChart elementId="s3BarChart" data={this.props.s3Data}/>}
            </div>
          </div>

        </div>

        <div className="white-box no-padding">
          {this.props.s3Data && <Components.S3Analytics.S3AnalyticsTable data={this.props.s3Data}/>}
        </div>
      </div>
    );
  }

}

S3AnalyticsContainer.propTypes = {};

const mapStateToProps = ({aws}) => ({s3Data: aws.s3.data});

const mapDispatchToProps = (dispatch) => ({
  getS3Data: () => {
    dispatch(Actions.AWS.S3.getS3Data())
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(S3AnalyticsContainer);
