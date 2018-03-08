import React, { Component } from 'react';
import PropTypes from 'prop-types';
import NVD3Chart from 'react-nvd3';
import {s3Analytics} from '../../../common/formatters';
import Spinner from 'react-spinkit';
import 'nvd3/build/nv.d3.min.css';
import * as d3 from "d3";

const transformStoragePieChart = s3Analytics.transformStoragePieChart;
const transformBandwidthPieChart = s3Analytics.transformBandwidthPieChart;
const transformRequestsPieChart = s3Analytics.transformRequestsPieChart;
const getTotalPieChart = s3Analytics.getTotalPieChart;

/* istanbul ignore next */
const formatX = (d) => (d.key);

/* istanbul ignore next */
const formatY = (d) => (d.value);

// PieChart Component
export class PieChartComponent extends Component {

  generateDatum = () => {
    if (Object.keys(this.props.data.values).length)
      switch (this.props.mode) {
        case "storage":
          return transformStoragePieChart(this.props.data.values);
        case "bandwidth":
          return transformBandwidthPieChart(this.props.data.values);
        case "requests":
          return transformRequestsPieChart(this.props.data.values);
        default:
          return null;
      }
    return null;
  };

  render() {
    if (!this.props.data || !this.props.data.status)
      return (<Spinner className="spinner clearfix" name='circle'/>);

    if (this.props.data && this.props.data.status && this.props.data.hasOwnProperty("error"))
      return (<div className="alert alert-warning" role="alert">Data not available ({this.props.data.error.message})</div>);

    const datum = this.generateDatum();
    if (!datum)
      return (<h3 className="no-data">No data available</h3>);
    const total = '$' + d3.format(',.2f')(getTotalPieChart(datum));
    return (
      <div className="s3analytics piechart">
        <h2>{this.props.mode} Cost</h2>
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
          height={270}
        />
      </div>
    )
  }

}

PieChartComponent.propTypes = {
  data: PropTypes.object,
  mode: PropTypes.oneOf(["storage", "bandwidth", "requests"])
};

class StorageCostChartComponent extends Component {

  render() {
    return (<PieChartComponent mode="storage" data={this.props.data}/>)
  }

}

class BandwidthCostChartComponent extends Component {

  render() {
    return (<PieChartComponent mode="bandwidth" data={this.props.data}/>)
  }

}

class RequestsCostChartComponent extends Component {

  render() {
    return (<PieChartComponent mode="requests" data={this.props.data}/>)
  }

}

export default {
  StorageCostChartComponent,
  BandwidthCostChartComponent,
  RequestsCostChartComponent
};
