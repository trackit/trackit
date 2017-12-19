import React, { Component } from 'react';
import { connect } from 'react-redux';
import NVD3Chart from 'react-nvd3';
import * as d3 from 'd3';
import Components from "../../components";

import Actions from "../../actions";

import 'nvd3/build/nv.d3.min.css';

const TimerangeSelector = Components.Misc.TimerangeSelector;
const Panel = Components.Misc.Panel;

class CostBreakdownContainer extends Component {

  componentWillMount() {
    this.props.getCosts(this.props.costsDates.startDate, this.props.costsDates.endDate, ["product", this.props.costsInterval]);
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.costsDates !== nextProps.costsDates || this.props.costsInterval !== nextProps.costsInterval)
      nextProps.getCosts(nextProps.costsDates.startDate, nextProps.costsDates.endDate, ["product", nextProps.costsInterval]);
  }

  transformProducts = (data, timescope) => {
    try {
      let dates = [];
      Object.keys(data.product).forEach((key) => {
        Object.keys(data.product[key][timescope]).forEach((date) => {
          if (dates.indexOf(date) === -1)
            dates.push(date);
        })
      });
      return Object.keys(data.product).map((key) => ({
        key,
        values: dates.map((date) => ([date, data.product[key][timescope][date] || 0]))
      }));
    } catch (e) {
      return [];
    }
  };

  render() {

    const context = {
      formatXAxis: (d) => (d3.time.format('%x')(new Date(d))),
      formatYAxis: (d) => ('$' + d3.format(',.2f')(d)),
    };

    const datum = (this.props.costsValues && this.props.costsInterval ? this.transformProducts(this.props.costsValues, this.props.costsInterval) : null);

    const chart = (this.props.costsValues && this.props.costsInterval ? (
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
    ) : null);

    return(
      <Panel>

        <div>
          <h3 className="white-box-title no-padding inline-block">
            Cost Breakdown
          </h3>
          <div className="inline-block pull-right">
            <TimerangeSelector
              startDate={this.props.costsDates.startDate}
              endDate={this.props.costsDates.endDate}
              setDatesFunc={this.props.setCostsDates}
              interval={this.props.costsInterval}
              setIntervalFunc={this.props.setCostsInterval}
            />
          </div>
        </div>

        <div>
          {chart}
        </div>

      </Panel>
    );
  }
}

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  costsValues: aws.costs.values,
  costsDates: aws.costs.dates,
  costsInterval: aws.costs.interval
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getCosts: (begin, end, filters, accounts=undefined) => {
    dispatch(Actions.AWS.Costs.getCosts(begin, end, filters, accounts));
  },
  setCostsDates: (startDate, endDate) => {
    dispatch(Actions.AWS.Costs.setCostsDates(startDate, endDate))
  },
  setCostsInterval: (interval) => {
    dispatch(Actions.AWS.Costs.setCostsInterval(interval));
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(CostBreakdownContainer);
