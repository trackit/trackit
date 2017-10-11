import React, { Component } from 'react';
import {bindActionCreators} from 'redux';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { CircularProgress } from 'material-ui/Progress';
import Switch from 'material-ui/Switch';
import { FormControlLabel } from 'material-ui/Form';


import { dataToBarChart } from '../common/formatters';
import BarChartComponent from './BarChartComponent';
import * as ProvidersActions from '../actions/providersActions';


const googleColor = '#4885ed';
const awsColor = '#ff9900';

// ChartManagerComponent Component
class ChartManagerComponent extends Component {

    componentDidMount() {}

    render() {

      const styles = {
        spinner: {
          margin: '50px auto',
          width: '50px',
          height: '50px',
          display: 'block',
        },
      };

      let gcpChart;
      if (this.props.gcp.pricing) {
        const gcpData = dataToBarChart(this.props.gcp.pricing, googleColor);
        gcpChart = <BarChartComponent
          elementId="gcpChart"
          data={gcpData}
          barmode='group'
          title='Google Cloud Platform Storage (monthly price)'
        />;
      } else {
        gcpChart = <CircularProgress style={styles.spinner} />;
      }

      let awsChart;
      if (this.props.aws.pricing) {
        const awsData = dataToBarChart(this.props.aws.pricing, awsColor);
        awsChart = <BarChartComponent
          elementId="awsChart"
          data={awsData}
          barmode='group'
          title='AWS S3 Storage (monthly price)'
        />;
      } else {
        awsChart = <CircularProgress style={styles.spinner} />;
      }


      return (
        <div>
          <div>
            <div className="col-md-12">
              <h4 className="paper-title">
                <i className="fa fa-bar-chart red-color"/>
                &nbsp;
                Charts
              </h4>
              <FormControlLabel
                className="pull-right"
                control={
                  <Switch
                    id="stackedViewSwitch"
                    checked={false}
                    onChange={() => {}}
                    aria-label="Chart grouped view switch"
                  />
                }
                label="Switch to stacked view ?"
              />
            </div>
          </div>
          <div className="clearfix" />
          <hr />
          <div className="row">
            <div className="col-md-12">
              {gcpChart}
            </div>
          </div>
          <div className="row">
            <div className="col-md-12">
              {awsChart}
            </div>
          </div>
        </div>
      );
    }

}


// Define PropTypes
ChartManagerComponent.propTypes = {
  gcp: PropTypes.object,
  aws: PropTypes.object,
};


// Subscribe component to redux store and merge the state into
// component's props
const mapStateToProps = ({ gcp, aws }) => ({
  gcp,
  aws
});

const mapActionCreatorsToProps = (dispatch) => (
   bindActionCreators(ProvidersActions, dispatch)
);


// connect method from react-router connects the component with redux store
export default connect(mapStateToProps, mapActionCreatorsToProps)(ChartManagerComponent);
