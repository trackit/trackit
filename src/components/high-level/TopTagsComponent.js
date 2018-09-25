import React, {Component} from 'react';
import PropTypes from 'prop-types';
import moment from 'moment';
import NVD3Chart from 'react-nvd3';
import * as d3 from 'd3';
import {formatChartPrice} from '../../common/formatters';


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

class TopTags extends Component {

    handleKeySelection(event) {
        this.props.setSelected(event.target.value);
    }
    
    findProduct(toFind, values) {
        let result = null;
        values.forEach((product) => {
        if (product.key === toFind)
            result = product;
        });
        return result;
    }

    getProjectedValues(values) {
        return values.map((item) => {
            if (item.costs[1] && item.costs[1].cost) {
                return ([item.tag.length ? item.tag : 'No value', (item.costs[1].cost / moment().date()) * parseInt(moment().endOf('month').format("DD"), 10)]);
            } else {
                return ([item.tag.length ? item.tag : 'No value', (item.costs[0].cost / moment().date()) * parseInt(moment().endOf('month').format("DD"), 10)]);
            }
        });
    }

    formatDataForChart(values) {
        const res = [
            {
                key: this.props.currentInterval ? 'Current' : 'Selected month',
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
            const tag = element.tag.length ? element.tag : 'No value';
            if (element.costs[0] && element.costs[0].cost && element.costs[1] && element.costs[1].cost) {
                res[0].values.push([tag, element.costs[1].cost]);
                res[1].values.push([tag, element.costs[0].cost]);    
            } else if (element.costs[0] && element.costs[0].cost) {
                res[0].values.push([tag, element.costs[0].cost]);
                res[1].values.push([tag, 0]);    
            }
        }
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
        if (this.props.costs && this.props.costs.status && this.props.costs.values) {
            const values = JSON.parse(JSON.stringify(this.props.costs.values));

            if (!values.length)
              chart = (<h4 className="no-data">No data available for this timerange</h4>);
            else {
                const sorted = values.sort((a, b) =>  {
                    if (a.costs[1] && b.costs[1]) {
                        return (a.costs[1].cost > b.costs[1].cost ? -1 : (a.costs[1].cost < b.costs[1].cost ? 1 : 0));
                    } else if (a.costs[0] && b.costs[0]) {
                        return (a.costs[0].cost > b.costs[0].cost ? -1 : (a.costs[0].cost < b.costs[0].cost ? 1 : 0));
                    } else {
                        return 1;
                    }
                }).slice(0, 5);
                const datum = this.formatDataForChart(sorted);
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
            }    
        }

        let selector;
        if (this.props.selected) {
            selector = (
                <select className="hl-panel-select" style={{maxWidth: '110px'}} value={this.props.selected} onChange={this.handleKeySelection.bind(this)}>
                    {this.props.keys.map(item => <option key={item} value={item}>{item}</option>)}
                </select>
            );
        }


        return (
        <div className="col-md-6">
            <div className="white-box hl-panel">
                <h4 className="m-t-0 hl-panel-title">{moment(this.props.date).format('MMM Y')} Top 5 categories</h4>
                {selector}
                {chart}
            <div className="clearfix"/>
            </div>
        </div>
        );
    }
}

TopTags.propTypes = {
    date: PropTypes.object.isRequired,
    keys: PropTypes.array.isRequired,
    selected: PropTypes.string,
    currentInterval: PropTypes.bool.isRequired,
    setSelected: PropTypes.func.isRequired,
    costs: PropTypes.object.isRequired,
};

export default TopTags;