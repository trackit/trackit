import React, { Component } from 'react';
import PropTypes from 'prop-types';
import * as d3 from 'd3';
import NVD3Chart from 'react-nvd3';
import { costBreakdown } from '../../../common/formatters';
import 'nvd3/build/nv.d3.min.css';
import ChartsColors from "../../../styles/ChartsColors";

const transformProductsBarChart = costBreakdown.transformProductsBarChart;

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
  right: 10,
  left: 70,
};

class BarChartComponent extends Component {

  generateDatum = () => {
    if (this.props.values && Object.keys(this.props.values).length && this.props.interval && this.props.filter)
      return transformProductsBarChart(this.props.values, this.props.filter, this.props.interval);
    return null;
  };

  render() {
    const datum = this.generateDatum();
    if (!datum)
      return (<h4 className="no-data">No data available for this timerange</h4>);
    return (
      <div>
        <NVD3Chart
          id="barChart"
          type="multiBarChart"
          datum={datum}
          color={ChartsColors}
          context={context}
          xAxis={xAxis}
          yAxis={yAxis}
          margin={this.props.margin ? margin : null}
          rightAlignYAxis={false}
          clipEdge={false}
          showControls={true}
          showLegend={this.props.legend}
          stacked={true}
          x={formatX}
          y={formatY}
          height={this.props.height}
        />
      </div>
    )
  }

}

BarChartComponent.propTypes = {
  values: PropTypes.object,
  interval: PropTypes.string.isRequired,
  filter: PropTypes.string.isRequired,
  legend: PropTypes.bool.isRequired,
  height: PropTypes.number.isRequired,
  margin: PropTypes.bool
};

BarChartComponent.defaultProps = {
  margin: true
};

export default BarChartComponent;
