import React, { Component } from 'react';
import PropTypes from 'prop-types';
import NVD3Chart from 'react-nvd3';
import ReactTable from 'react-table';
import {costBreakdown, formatPrice} from '../../../common/formatters';
import 'nvd3/build/nv.d3.min.css';
import * as d3 from "d3";
import ChartsColors from "../../../styles/ChartsColors";

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

  constructor(props) {
    super(props);
    this.state = {
      datum: [],
      total: 0
    };
    this.getSelectedTotal = this.getSelectedTotal.bind(this);
  }

  componentWillMount() {
    const datum = this.generateDatum();
    const total = '$' + d3.format(',.2f')(getTotalPieChart(datum));
    this.setState({datum, total});
  }

  generateDatum = () => {
    if (this.props.values && Object.keys(this.props.values).length && this.props.filter)
      return transformProductsPieChart(this.props.values, this.props.filter);
    return null;
  };

  getSelectedTotal = (selection, chart) => {
    const datum = [];
    this.state.datum.forEach((item, index) => {
      if (!selection[index])
        datum.push(item);
    });
    const total = '$' + d3.format(',.2f')(getTotalPieChart(datum));
    this.setState({total});
    chart.title(total);
    chart.update();
  };

  render() {
    if (!this.state.datum)
      return (<h4 className="no-data">No data available for this timerange</h4>);


    /* istanbul ignore next */
    const table = (this.props.table ? (
      <ReactTable
        data={this.state.datum}
        noDataText="No buckets available"
        columns={[
          {
            Header: 'Name',
            accessor: 'key',
            Cell: row => (<strong>{row.value}</strong>)
          }, {
            Header: 'Cost',
            accessor: 'value',
            Cell: row => (<span className="total-cell">{formatPrice(row.value)}</span>)
          }
        ]}
        defaultSorted={[{
          id: 'Cost',
          desc: true
        }]}
        defaultPageSize={10}
        className=" -highlight"
      />
    ) : null);

    const chart = (
      <NVD3Chart
        id="pieChart"
        type="pieChart"
        title={this.state.total}
        datum={this.state.datum}
        color={ChartsColors}
        margin={this.props.margin ? margin : null}
        x={formatX}
        y={formatY}
        showLabels={false}
        showLegend={this.props.legend}
        legendPosition="right"
        donut={true}
        height={this.props.height}
        renderStart={(chart) => {
          chart.dispatch.on('stateChange', (data) => {
            this.getSelectedTotal(data.disabled, chart);
          })
        }}
      />
    );

    return (
      <div className="clearfix">
        {chart}
        {table}
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
  margin: PropTypes.bool,
  table: PropTypes.bool
};

PieChartComponent.defaultProps = {
  margin: true,
  table: false
};

export default PieChartComponent;