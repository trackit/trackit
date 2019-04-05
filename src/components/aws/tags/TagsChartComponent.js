import React, { Component } from 'react';
import PropTypes from 'prop-types';
import NVD3Chart from 'react-nvd3';
import ReactTable from 'react-table';
import {tags, formatPrice} from '../../../common/formatters';
import 'nvd3/build/nv.d3.min.css';
import * as d3 from "d3";
import ChartsColors from "../../../styles/ChartsColors";

const transformItemsPieChart = tags.transformItemsPieChart;
const getTotalPieChart = tags.getTotalPieChart;

/* istanbul ignore next */
const formatX = (d) => (d.key);

/* istanbul ignore next */
const formatY = (d) => (d.value);

const margin = {
  right: 0
};

class TagsChartComponent extends Component {

  constructor(props) {
    super(props);
    this.state = {
      datum: [],
      total: 0
    };
    this.getSelectedTotal = this.getSelectedTotal.bind(this);
  }

  componentWillMount() {
    const datum = transformItemsPieChart(this.props.values);
    const total = '$' + d3.format(',.2f')(getTotalPieChart(datum));
    this.setState({datum, total});
  }

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

    const itemsList = [];
    this.state.datum.forEach((tag) => {
      Object.keys(tag.items).forEach((item) => {
        if (itemsList.indexOf(item) === -1)
          itemsList.push(item);
      })
    });
    const itemsColumns = itemsList.map((item) => ({
      Header: (item && item.length ? item : `No ${this.props.filter}`),
      accessor: 'items',
      id: item,
      sortMethod: (a, b) => (a.hasOwnProperty(item) && b.hasOwnProperty(item) && a[item] > b[item] ? 1 : -1),
      Cell: row => (<span className="total-cell">{formatPrice(row.value[item] || 0)}</span>)
    }));

    /* istanbul ignore next */
    const table = (
      <ReactTable
        data={this.state.datum}
        noDataText="No tags available"
        columns={[
          {
            Header: (<strong>Tag</strong>),
            accessor: 'key',
            Cell: row => (<strong>{row.value}</strong>)
          }, {
            Header: (<strong>Total</strong>),
            accessor: 'value',
            Cell: row => (<strong className="total-cell">{formatPrice(row.value)}</strong>)
          }, {
            Header: `${this.props.filter}s`,
            columns: itemsColumns
          }
        ]}
        defaultSorted={[{
          id: 'Cost',
          desc: true
        }]}
        defaultPageSize={10}
        className=" -highlight"
      />
    );

    return (
      <div className="m-t-20">
        <div className="row">
          <div className="col-md-3">
            <NVD3Chart
              id="pieChart"
              type="pieChart"
              title={this.state.total}
              datum={this.state.datum}
              color={ChartsColors}
              margin={margin}
              x={formatX}
              y={formatY}
              showLabels={false}
              showLegend={this.props.legend}
              legendPosition="top"
              donut={true}
              height={this.props.height}
              renderStart={(chart) => {
                chart.dispatch.on('stateChange', (data) => {
                  this.getSelectedTotal(data.disabled, chart);
                })
              }}
            />
          </div>
          <div className="col-md-9">
            {table}
          </div>
        </div>
      </div>
    )
  }

}

TagsChartComponent.propTypes = {
  values: PropTypes.arrayOf(PropTypes.object),
  legend: PropTypes.bool.isRequired,
  height: PropTypes.number.isRequired,
  filter: PropTypes.string.isRequired
};


export default TagsChartComponent;
