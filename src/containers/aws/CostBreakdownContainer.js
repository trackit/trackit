import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import NVD3Chart from 'react-nvd3';
import * as d3 from 'd3';
import Components from "../../components";

import Actions from "../../actions";

import s3square from '../../assets/s3-square.png';
import 'nvd3/build/nv.d3.min.css';

const TimerangeSelector = Components.Misc.TimerangeSelector;
const Selector = Components.Misc.Selector;
const Panel = Components.Misc.Panel;

export class CostBreakdownContainer extends Component {

  componentWillMount() {
    this.props.getCosts(this.props.costsDates.startDate, this.props.costsDates.endDate, [this.props.costsFilter, this.props.costsInterval]);
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.costsDates !== nextProps.costsDates ||
      this.props.costsInterval !== nextProps.costsInterval ||
      this.props.costsFilter !== nextProps.costsFilter)
      nextProps.getCosts(nextProps.costsDates.startDate, nextProps.costsDates.endDate, [nextProps.costsFilter, nextProps.costsInterval]);
  }

  transformProducts = (data, filter, interval) => {
    if (filter === "all" && data.hasOwnProperty(interval))
      return [{
        key: "Total",
        values: Object.keys(data[interval]).map((date) => ([date, data[interval][date]]))
      }];
    else if (!data.hasOwnProperty(filter))
      return [];
    let dates = [];
    try {
      Object.keys(data[filter]).forEach((key) => {
        Object.keys(data[filter][key][interval]).forEach((date) => {
          if (dates.indexOf(date) === -1)
            dates.push(date);
        })
      });
      return Object.keys(data[filter]).map((key) => ({
        key: (key.length ? key : `No ${filter}`),
        values: dates.map((date) => ([date, data[filter][key][interval][date] || 0]))
      }));
    } catch (e) {
      return [];
    }
  };

  render() {

    const filters = {
      all: "Total",
      account: "Account",
      product: "Product",
      region: "Region"
    };

    /* istanbul ignore next */
    const context = {
      formatXAxis: (d) => (d3.time.format('%x')(new Date(d))),
      formatYAxis: (d) => ('$' + d3.format(',.2f')(d)),
    };

    const datum = (this.props.costsValues && this.props.costsInterval && this.props.costsFilter ? this.transformProducts(this.props.costsValues, this.props.costsFilter, this.props.costsInterval) : null);

    const chart = (datum ? (
      <NVD3Chart
        id="barChart"
        type="multiBarChart"
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
        margin={{right:100}}
        rightAlignYAxis={true}
        clipEdge={true}
        showControls={true}
        x={
          /* istanbul ignore next */
          (d) => {
          const date = new Date(d[0]);
          return date.getTime();
        }}
        y={
          /* istanbul ignore next */
          (d) => (d[1])
        }
        height={600}
      />
    ) : null);

    return(
      <Panel>

        <div>
          <h3 className="white-box-title no-padding inline-block">
            <img className="white-box-title-icon" src={s3square} alt="AWS square logo"/>
            Cost Breakdown
          </h3>

          <div className="inline-block pull-right">
            <div className="inline-block">
              <Selector
                values={filters}
                selected={this.props.costsFilter}
                selectValue={this.props.setCostsFilter}
              />
            </div>
            <div className="inline-block">
              <TimerangeSelector
                startDate={this.props.costsDates.startDate}
                endDate={this.props.costsDates.endDate}
                setDatesFunc={this.props.setCostsDates}
                interval={this.props.costsInterval}
                setIntervalFunc={this.props.setCostsInterval}
              />
            </div>
          </div>

          <div className="clearfix"/>

        </div>

        <div>
          {chart}
        </div>

      </Panel>
    );
  }
}

CostBreakdownContainer.propTypes = {
  costsValues: PropTypes.object,
  costsDates: PropTypes.shape({
    startDate: PropTypes.object,
    endDate: PropTypes.object,
  }),
  costsInterval:PropTypes.string.isRequired,
  costsFilter:PropTypes.string.isRequired,
  getCosts: PropTypes.func.isRequired,
  setCostsDates: PropTypes.func.isRequired,
  setCostsInterval: PropTypes.func.isRequired,
  setCostsFilter: PropTypes.func.isRequired,
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  costsValues: aws.costs.values,
  costsDates: aws.costs.dates,
  costsInterval: aws.costs.interval,
  costsFilter: aws.costs.filter
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getCosts: (begin, end, filters) => {
    dispatch(Actions.AWS.Costs.getCosts(begin, end, filters));
  },
  setCostsDates: (startDate, endDate) => {
    dispatch(Actions.AWS.Costs.setCostsDates(startDate, endDate))
  },
  setCostsInterval: (interval) => {
    dispatch(Actions.AWS.Costs.setCostsInterval(interval));
  },
  setCostsFilter: (filter) => {
    dispatch(Actions.AWS.Costs.setCostsFilter(filter));
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(CostBreakdownContainer);
