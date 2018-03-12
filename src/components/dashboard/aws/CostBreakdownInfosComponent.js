import React, { Component } from 'react';
import PropTypes from 'prop-types';
import AWS from '../../aws';

const Infos = AWS.CostBreakdown.Infos;

class CostBreakdownInfosComponent extends Component {

  constructor(props) {
    super(props);
    this.getValues = this.getValues.bind(this);
  }

  getValues(id, startDate, endDate, filters) {
    this.props.getValues(id, "costbreakdown", startDate, endDate, filters);
  }

  render() {
    return (
      <Infos
        type="bar"
        legend={false}
        height={500}
        id={this.props.id}
        accounts={this.props.accounts}
        values={this.props.values}
        getCosts={this.getValues}
        dates={this.props.dates}
        setDates={this.props.setDates}
        interval={this.props.interval}
        setInterval={this.props.setInterval}
      />
    );
  }

}

CostBreakdownInfosComponent.propTypes = {
  id: PropTypes.string.isRequired,
  accounts: PropTypes.arrayOf(PropTypes.object),
  values: PropTypes.object,
  getValues: PropTypes.func.isRequired,
  dates: PropTypes.shape({
    startDate: PropTypes.object,
    endDate: PropTypes.object,
  }),
  setDates: PropTypes.func.isRequired,
  interval: PropTypes.string.isRequired,
  setInterval: PropTypes.func.isRequired
};

export default CostBreakdownInfosComponent;
