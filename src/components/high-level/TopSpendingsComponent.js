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
        for (let i = 0; i < values.length; i++) {
            const element = values[i];
            if (element.key === toFind) {
                return element;
            }
        }
        return null;
    }

    getProjectedValues(values) {
        const res = [];
        for (let i = 0; i < values.length; i++) {
            const element = values[i];
            const projectedValue = (element.value / moment().date()) * parseInt(moment().endOf('month').format("DD"), 10);
          res.push([element.key, projectedValue]);
        }
        return res;
    }

    formatDataForChart(values, previousValues, isCurrent) {
        const res = [
            {
                key: isCurrent ? 'Current' : 'Selected month',
                values: [],
                color: '#0088CC'
            },
            {
                key: 'Previous Month',
                values: [],
                color: '#62bed2'
            },
        ];
        for (let i = 0; i < values.length; i++) {
            const element = values[i];
            const previousElement = previousValues[i];
            res[0].values.push([element.key, element.value]);
            if (previousElement && previousElement.key && previousElement.value)
                res[1].values.push([previousElement.key, previousElement.value]);
        }
        if (isCurrent) {
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
        const months = Object.keys(this.props.costs.months);
        const isSelectedCurrent = moment(months[1]).month() === moment().month();
        let currentMonthProducts = convertProductObjectToArray(this.props.costs.months[months[1]].product);
        const previousProducts = convertProductObjectToArray(this.props.costs.months[months[0]].product);
        // Sorting by price
        currentMonthProducts.sort((a, b) => {
            if (a.value > b.value) {
                return -1;
            } else if (a.value < b.value) {
                return 1;
            }
            return 0;
        });
        if (currentMonthProducts.length >= 5) {
            currentMonthProducts = currentMonthProducts.slice(0, 5);
        }
        const mappedPreviousMonthProducts = [];
        for (let i = 0; i < currentMonthProducts.length; i++) {
            const element = currentMonthProducts[i];
            mappedPreviousMonthProducts.push(this.findProduct(element.key, previousProducts));
        }
        const datum = this.formatDataForChart(currentMonthProducts, mappedPreviousMonthProducts, isSelectedCurrent);
        return (
            <div className="col-md-6">
                <div className="white-box hl-panel">
                    <h4 className="m-t-0 hl-panel-title">
                        {moment(this.props.date).format('MMM Y')} Top 5 spendings
                    </h4>
                    <Link to="/app/costbreakdown" className="hl-details-link">
                        More details
                    </Link>
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
                    <div className="clearfix"></div>
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
};

export default TopSpendings;