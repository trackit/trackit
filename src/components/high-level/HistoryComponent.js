import React, { Component } from 'react';
import PropTypes from 'prop-types';
import NVD3Chart from 'react-nvd3';
import * as d3 from 'd3';
import { formatChartPrice } from '../../common/formatters';

const convertProductObjectToArray = (values) => Object.keys(values).map((key) => ({key, value: values[key]}));

/* istanbul ignore next */
const context = {
  formatXAxis: (d) => (d3.time.format('%m/%Y')(new Date(d))),
  formatYAxis: (d) => (formatChartPrice(d)),
  valueFormatter: d => (`$${d3.format(',.0f')(d)}`),
};

const xAxis = {
  tickFormat: {
    name: 'formatXAxis',
    type: 'function',
  }
};

const yAxis = {
  tickFormat: {
    name: 'formatYAxis',
    type: 'function',
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
  right: 30,
  left: 25,
  bottom: 20,
};


class HistoryComponent extends Component {

  formatDataForChart(values) {
    const res = [
      {
        key: 'Cost',
        values: [],
        color: '#4885ed',
        area: true,
      },
    ];
    for (let i = 0; i < values.length; i++) {
      const element = values[i];
      res[0].values.push([element.key, element.value]);
    }
    return res;
  }

  getMax(values) {
    let max = 0;
    for (const key in values) {
      if (values.hasOwnProperty(key)) {
        const element = values[key];
        if (element > max) {
          max = element;
        }
      }
    }
    return max;
  }

  render() {
    let historyData = convertProductObjectToArray(this.props.history);
    const datum = this.formatDataForChart(historyData);
    return (
      <div className="col-md-12">
        <div className="white-box">
          <h4 className="m-t-0 hl-panel-title">History (Last 12 months)</h4>
          <NVD3Chart
            id="lineChart"
            type="lineChart"
            datum={datum}
            context={context}
            xAxis={xAxis}
            yAxis={yAxis}
            margin={margin}
            rightAlignYAxis={true}
            clipEdge={false}
            showControls={false}
            x={formatX}
            y={formatY}
            height={200}
            interpolate={'monotone'}
            useInteractiveGuideline={true}
            yDomain={[0, this.getMax(this.props.history)]}
            interactiveLayer={{
              tooltip: {
                valueFormatter: {
                  name: 'valueFormatter',
                  type: 'function',
                },
              }
            }}
          />
          <div className="clearfix"></div>
        </div>
      </div>
    );
  }
}

HistoryComponent.propTypes = {
  history: PropTypes.object.isRequired,
}

export default HistoryComponent;