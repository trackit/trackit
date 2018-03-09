import React, { Component } from 'react';
import PropTypes from 'prop-types';
import NVD3Chart from 'react-nvd3';
import {costBreakdown} from '../../../common/formatters';
import 'nvd3/build/nv.d3.min.css';
import * as d3 from "d3";

const transformProductsPieChart = costBreakdown.transformProductsPieChart;
const getTotalPieChart = costBreakdown.getTotalPieChart;

/* istanbul ignore next */
const formatX = (d) => (d.key);

/* istanbul ignore next */
const formatY = (d) => (d.value);

const margin = {
  right: 100
};

class PieChartComponent extends Component {

  generateDatum = () => {
    if (this.props.values && Object.keys(this.props.values).length && this.props.filter)
      return transformProductsPieChart(this.props.values, this.props.filter);
    return null;
  };

  render() {
    const datum = this.generateDatum();
    if (!datum)
      return (<h4 className="no-data">No data available for this timerange</h4>);
    const total = '$' + d3.format(',.2f')(getTotalPieChart(datum));
    return (
      <div>
        <NVD3Chart
          id="pieChart"
          type="pieChart"
          title={total}
          datum={datum}
          margin={this.props.margin ? margin : null}
          x={formatX}
          y={formatY}
          showLabels={false}
          showLegend={this.props.legend}
          legendPosition="right"
          donut={true}
          height={(this.props.values && Object.keys(this.props.values).length ? this.props.height : 150)}
        />
      </div>
    )
  }

}

PieChartComponent.propTypes = {
  values: PropTypes.object,
  interval: PropTypes.string.isRequired,
  filter: PropTypes.string.isRequired,
  legend: PropTypes.bool.isRequired,
  height: PropTypes.number.isRequired,
  margin: PropTypes.bool
};

PieChartComponent.defaultProps = {
  margin: true
};

export default PieChartComponent;