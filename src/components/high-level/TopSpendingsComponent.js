import React, {Component} from 'react';
import PropTypes from 'prop-types';
import moment from 'moment';
import {Link} from 'react-router-dom';
import NVD3Chart from 'react-nvd3';
import * as d3 from 'd3';
import {formatChartPrice} from '../../common/formatters';

const convertProductObjectToArray = (values) => {
  const res = Object.keys(values).map((key) => {
    return {key, value: values[key]}
  });
  return res;
};

/* istanbul ignore next */
const context = {
  formatYAxis: (d) => (formatChartPrice(d)),
  valueFormatter: d => (`$${d3.format(',.0f')(d)}`),
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
  right: 70,
  left: 10,
  bottom: 80,
};

class TopSpendings extends Component {

  findProduct(toFind, values) {
    let result = null;
    values.forEach((product) => {
      if (product.key === toFind)
        result = product;
    });
    return result;
  }

  getProjectedValues(values) {
    return values.map((product) => [product.key, (product.value / moment().date()) * parseInt(moment().endOf('month').format("DD"), 10)]);
  }

  formatDataForChart(values, previousValues) {
    const res = [
      {
        key: this.props.currentInterval ? 'Current' : 'Selected month',
        values: values.length ? values.map((product) => (product !== null ? [product.key, product.value] : [])) : [],
        color: '#0088CC'
      },
      {
        key: 'Previous Month',
        values: previousValues.length ? previousValues.map((product) => (product !== null ? [product.key, product.value] : [])) : [],
        color: '#62bed2'
      },
    ];
    if (this.props.currentInterval) {
      res[2] = res[1];
      res[1] = {
        key: 'Projected cost',
        values: this.getProjectedValues(values),
        color: '#FFBF00',
      };
    }
    return res;
  }

  render() {
    let chart;

    const months = Object.keys(this.props.costs.months);

    if (!months.length)
      chart = (<h4 className="no-data">No data available for this timerange</h4>);
    else {
      let selectedMonthProducts = [];
      let previousMonthProducts = [];
      let mappedPreviousMonthProducts = [];
      const parsedMonths = months.map((month) => (moment(month)));

      if (parsedMonths.length === 2) {
        selectedMonthProducts = convertProductObjectToArray(this.props.costs.months[months[1]].product);
        previousMonthProducts = convertProductObjectToArray(this.props.costs.months[months[0]].product);
      } else if (parsedMonths[0].isSame(this.props.date, "month")) {
        selectedMonthProducts = convertProductObjectToArray(this.props.costs.months[months[0]].product);
      }

      if (selectedMonthProducts.length) {
        // Sorting by price and limit to 5
        selectedMonthProducts = selectedMonthProducts.sort((a, b) => (a.value > b.value ? -1 : (a.value < b.value ? 1 : 0)))
                                                      .slice(0, 5);
        mappedPreviousMonthProducts = selectedMonthProducts.map((product) => (this.findProduct(product.key, previousMonthProducts)));
        const datum = this.formatDataForChart(selectedMonthProducts, mappedPreviousMonthProducts);
        chart = (
          <NVD3Chart
            id="barChart"
            type="multiBarChart"
            datum={datum}
            context={context}
            yAxis={yAxis}
            margin={margin}
            rightAlignYAxis={true}
            clipEdge={false}
            showControls={false}
            rotateLabels={25}
            reduceXTicks={false}
            x={formatX}
            y={formatY}
            height={350}
            tooltip={{
              valueFormatter: {
                name: 'valueFormatter',
                type: 'function',
              },
            }}
          />
        );
      } else
        chart = (<h4 className="no-data">No data available for this timerange</h4>);
    }

    return (
      <div className="col-md-6">
        <div className="white-box hl-panel">
          <h4 className="m-t-0 hl-panel-title">{moment(this.props.date).format('MMM Y')} Top 5 spendings</h4>
          <Link to="/app/costbreakdown" className="hl-details-link">
            More details
          </Link>
          {chart}
          <div className="clearfix"/>
        </div>
      </div>
    );
  }
}

TopSpendings.propTypes = {
  costs: PropTypes.shape({
    months : PropTypes.object.isRequired,
  }).isRequired,
  date: PropTypes.object.isRequired,
  currentInterval: PropTypes.bool.isRequired
};

export default TopSpendings;