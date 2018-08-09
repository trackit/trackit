import React, { Component } from 'react';
import PropTypes from 'prop-types';
import NVD3Chart from 'react-nvd3';
import ReactTable from 'react-table';
import {tags, formatPrice} from '../../../common/formatters';
import 'nvd3/build/nv.d3.min.css';
import * as d3 from "d3";
import ChartsColors from "../../../styles/ChartsColors";

const transformProductsPieChart = tags.transformProductsPieChart;
const getTotalPieChart = tags.getTotalPieChart;

/* istanbul ignore next */
const formatX = (d) => (d.key);

/* istanbul ignore next */
const formatY = (d) => (d.value);

const margin = {
  right: 0
};

class TagsChartComponent extends Component {

  render() {
    const datum = transformProductsPieChart(this.props.values);

    if (!datum)
      return (<h4 className="no-data">No data available for this timerange</h4>);

    const total = '$' + d3.format(',.2f')(getTotalPieChart(datum));

    const productsList = [];
    datum.forEach((tag) => {
      Object.keys(tag.products).forEach((product) => {
        if (productsList.indexOf(product) === -1)
          productsList.push(product);
      })
    });
    const productsColumns = productsList.map((product) => ({
      Header: product,
      accessor: 'products',
      id: product,
      sortMethod: (a, b) => (a.hasOwnProperty(product) && b.hasOwnProperty(product) && a[product] > b[product] ? 1 : -1),
      Cell: row => (<span className="total-cell">{formatPrice(row.value[product] || 0)}</span>)
    }));

    /* istanbul ignore next */
    const table = (
      <ReactTable
        data={datum}
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
            Header: 'Products',
            columns: productsColumns
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
              title={total}
              datum={datum}
              color={ChartsColors}
              margin={margin}
              x={formatX}
              y={formatY}
              showLabels={false}
              showLegend={this.props.legend}
              legendPosition="top"
              donut={true}
              height={this.props.height}
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
};


export default TagsChartComponent;
