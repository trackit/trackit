import React, { Component } from 'react';
import { connect } from 'react-redux';
import NVD3Chart from 'react-nvd3';
import * as d3 from 'd3';
import Panel from '../../components/misc/Panel';

import Actions from "../../actions";

class CostBreakdownContainer extends Component {

  componentWillMount() {
    this.props.getCosts("2017-10-01", "2017-11-01", ["product", "month"]);
  }

  transformProducts = (data, timeScope) => {
    let dates = [];
    Object.keys(data.product).forEach((key) => {
      Object.keys(data.product[key][timeScope]).forEach((date) => {
        if (dates.indexOf(date) === -1)
          dates.push(date);
      })
    });
    return Object.keys(data.product).map((key) => ({
      key,
      values: dates.map((date) => ([date, data.product[key][timeScope][date] || 0]))
    }));
  };

  render() {
    if (!this.props.costs)
      return null;

    const datum = this.transformProducts(this.props.costs, "month");

    console.log(datum);

    const context = {
      formatXAxis: (d) => (d3.time.format('%x')(new Date(d))),
      formatYAxis: (d) => ('$' + d3.format(',.2f')(d)),
    };

    return (
      <Panel title="Cost breakdown">
        <NVD3Chart
          id="barChart"
          type={this.props.horizontal ? 'multiBarHorizontalChart' : 'multiBarChart'}
          datum={datum}
          context={context}
          xAxis={{
              tickFormat: {
              name:'formatXAxis',
              type:'function',
              }
          }}
          yAxis={{
              tickFormat: {
              name:'formatYAxis',
              type:'function',
              }
          }}
          margin={this.props.horizontal ? {left: 100} : {right:100}}
          rightAlignYAxis={true}
          clipEdge={true}
          showControls={true}
          stacked={this.props.horizontal}
          x={(d) => {
            const date = new Date(d[0]);
            return date.getTime();
          }}
          y={(d) => (d[1])}
          height={this.props.horizontal ? 850 : 400}
        />
      </Panel>
    );
  }
}

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  costs: aws.costs
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getCosts: (begin, end, filters, accounts=undefined) => {
    dispatch(Actions.AWS.Costs.getCosts(begin, end, filters, accounts));
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(CostBreakdownContainer);
