import React, { Component } from 'react';
import PropTypes from 'prop-types';
import AWS from '../../aws';

const Chart = AWS.CostBreakdown.Chart;

class CostBreakdownPieChartComponent extends Component {

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
        type="pie"
        legend={false}
        height={350}
        margin={false}
        id={this.props.id}
        accounts={this.props.accounts}
        values={this.props.values}
        getCosts={this.getValues}
        dates={this.props.dates}
        interval={this.props.interval}
        setInterval={this.props.setInterval}
        filter={this.props.filter}
        setFilter={this.props.setFilter}
        table={false}
        icon={false}
      />
    );
  }

}

CostBreakdownPieChartComponent.propTypes = {
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

export default CostBreakdownPieChartComponent;
