import React, { Component } from 'react';
import PropTypes from 'prop-types';
import NVD3Chart from 'react-nvd3';
import * as d3 from 'd3';
import { transformProducts } from '../../../common/formatters';
import 'nvd3/build/nv.d3.min.css';

/* istanbul ignore next */
const context = {
  formatXAxis: (d) => (d3.time.format('%x')(new Date(d))),
  formatYAxis: (d) => ('$' + d3.format(',.2f')(d)),
};

const xAxis = {
  tickFormat: {
    name:'formatXAxis',
    type:'function',
  }
};

const yAxis = {
  tickFormat: {
    name:'formatYAxis',
    type:'function',
  }
};

/* istanbul ignore next */
const formatX = (d) => {
  const date = new Date(d[0]);
  return date.getTime();
};

/* istanbul ignore next */
const formatY = (d) => (d[1]);

const margin = {
  right: 100
};

class ChartComponent extends Component {

  generateDatum = () => {
    if (this.props.values && this.props.interval && this.props.filter)
      return transformProducts(this.props.values, this.props.filter, this.props.interval);
    return null;
  };

  render() {
    const datum = this.generateDatum();
    if (!datum)
      return null;
    return (
      <NVD3Chart
        id="barChart"
        type="multiBarChart"
        datum={datum}
        context={context}
        xAxis={xAxis}
        yAxis={yAxis}
        margin={margin}
        rightAlignYAxis={true}
        clipEdge={true}
        showControls={true}
        stacked={true}
        x={formatX}
        y={formatY}
        height={600}
      />
    )
  }

}

ChartComponent.propTypes = {
  values: PropTypes.object,
  interval: PropTypes.string.isRequired,
  filter: PropTypes.string.isRequired,
};

export default ChartComponent;