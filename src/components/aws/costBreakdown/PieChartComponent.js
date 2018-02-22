import React, { Component } from 'react';
import PropTypes from 'prop-types';
import NVD3Chart from 'react-nvd3';
import { transformProductsPieChart, getTotalPieChart } from '../../../common/formatters';
import 'nvd3/build/nv.d3.min.css';
import * as d3 from "d3";

/* istanbul ignore next */
const formatX = (d) => (d.key);

/* istanbul ignore next */
const formatY = (d) => (d.value);

const margin = {
  right: 100
};

class PieChartComponent extends Component {

  generateDatum = () => {
    if (this.props.values && this.props.filter)
      return transformProductsPieChart(this.props.values, this.props.filter);
    return null;
  };

  render() {
    const datum = this.generateDatum();
    if (!datum)
      return null;
    const total = '$' + d3.format(',.2f')(getTotalPieChart(datum));
    return (
      <NVD3Chart
        id="pieChart"
        type="pieChart"
        title={total}
        datum={datum}
        margin={margin}
        x={formatX}
        y={formatY}
        showLabels={false}
        showLegend={true}
        donut={true}
        height={(this.props.values && Object.keys(this.props.values).length ? 400 : 150)}
      />
    )
  }

}

PieChartComponent.propTypes = {
  values: PropTypes.object,
  interval: PropTypes.string.isRequired,
  filter: PropTypes.string.isRequired,
};

export default PieChartComponent;