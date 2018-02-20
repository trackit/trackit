import React, { Component } from 'react';
import PropTypes from 'prop-types';
import NVD3Chart from 'react-nvd3';
import * as d3 from 'd3';
import { transformBuckets } from '../../../common/formatters';
import 'nvd3/build/nv.d3.min.css';

/* istanbul ignore next */
const context = {
  formatXAxis: (d) => (d),
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
const formatX = (d) => (d[0]);

/* istanbul ignore next */
const formatY = (d) => (d[1]);

const margin = {
  right: 100
};

// S3AnalyticsBarChart Component
class BarChart extends Component {

  generateDatum = () => {
    if (this.props.data)
      return transformBuckets(this.props.data);
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
        height={400}
      />
    )
  }

}

BarChart.propTypes = {
  elementId: PropTypes.string.isRequired,
  data: PropTypes.object
};

export default BarChart;
