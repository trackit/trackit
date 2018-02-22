import React, { Component } from 'react';
import PropTypes from 'prop-types';
import NVD3Chart from 'react-nvd3';
import {s3Analytics} from '../../../common/formatters';
import Spinner from 'react-spinkit';
import 'nvd3/build/nv.d3.min.css';
import * as d3 from "d3";

const transformStoragePieChart = s3Analytics.transformStoragePieChart;
const getTotalPieChart = s3Analytics.getTotalPieChart;

/* istanbul ignore next */
const formatX = (d) => (d.key);

/* istanbul ignore next */
const formatY = (d) => (d.value);

// S3AnalyticsBarChart Component
class StorageCostChartComponent extends Component {

  generateDatum = () => {
    if (Object.keys(this.props.data.values).length)
      return transformStoragePieChart(this.props.data.values);
    return null;
  };

  render() {
    if (!this.props.data || !this.props.data.status)
      return (<Spinner className="spinner clearfix" name='circle'/>);

    if (this.props.data && this.props.data.status && this.props.data.hasOwnProperty("error"))
      return (<div className="alert alert-warning" role="alert">Data not available ({this.props.data.error.message})</div>);

    const datum = this.generateDatum();
    if (!datum)
      return null;
    const total = '$' + d3.format(',.2f')(getTotalPieChart(datum));
    return (
      <div className="s3analytics piechart">
        <h2>Storage Cost</h2>
        <NVD3Chart
          id="pieChart"
          type="pieChart"
          title={total}
          datum={datum}
          x={formatX}
          y={formatY}
          showLabels={false}
          showLegend={false}
          donut={true}
          height={400}
        />
      </div>
    )
  }

}

StorageCostChartComponent.propTypes = {
  data: PropTypes.object
};

export default StorageCostChartComponent;
