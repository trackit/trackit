import React, { Component } from 'react';
import PropTypes from 'prop-types';
import AWS from '../../aws';

const Chart = AWS.CostBreakdown.Chart;

class CostBreakdownBarChartComponent extends Component {

  constructor(props) {
    super(props);
    this.getValues = this.getValues.bind(this);
  }

  getValues(id, startDate, endDate, filters) {
    this.props.getValues(id, "costbreakdown", startDate, endDate, filters);
  }

  render() {
    return (
      <Chart
        type="bar"
        legend={false}
        height={350}
        id={this.props.id}
        accounts={this.props.accounts}
        values={this.props.values}
        getCosts={this.getValues}
        dates={this.props.dates}
        interval={this.props.interval}
        setInterval={this.props.setInterval}
        filter={this.props.filter}
        setFilter={this.props.setFilter}
      />
    );
  }

}

CostBreakdownBarChartComponent.propTypes = {
  id: PropTypes.string.isRequired,
  accounts: PropTypes.arrayOf(PropTypes.object),
  values: PropTypes.object,
  getValues: PropTypes.func.isRequired,
  dates: PropTypes.shape({
    startDate: PropTypes.object,
    endDate: PropTypes.object,
  }),
  filter: PropTypes.string.isRequired,
  setFilter: PropTypes.func.isRequired,
  interval: PropTypes.string.isRequired,
  setInterval: PropTypes.func.isRequired
};

export default CostBreakdownBarChartComponent;
