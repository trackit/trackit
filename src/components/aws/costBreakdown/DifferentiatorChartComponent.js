import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ReactTable from 'react-table';
import Moment from 'moment';
import {costBreakdown, formatPrice, formatPercent, formatDate} from '../../../common/formatters';

const transformCostDifferentiator = costBreakdown.transformCostDifferentiator;

class DifferentiatorChartComponent extends Component {

  generateDatum = () => {
    if (this.props.values && Object.keys(this.props.values).length)
      return transformCostDifferentiator(this.props.values);
    return null;
  };

  /* istanbul ignore next */
  generateColumns(dates) {
    return dates.map((date, index) => {
      let columns = [{
        Header: 'Cost',
        id: date + '.cost',
        accessor: row => row[date].cost,
        Cell: row => (<span className="cost-cell">{formatPrice(row.value)}</span>)
      }];
      if (index > 0)
        columns.push({
          Header: 'Variation',
          id: date + '.variation',
          accessor: row => row[date].variation,
          Cell: row => (<span className="percentvariation-cell">{formatPercent(row.value)}</span>)
        });
      return ({
        Header: formatDate(Moment(date), this.props.interval),
        columns
      })
    })
  }

  render() {
    const datum = this.generateDatum();

    if (!datum)
      return (<h4 className="no-data">No data available for this timerange</h4>);

    const dates = this.generateColumns(datum.dates);

    /* istanbul ignore next */
    return (
      <div className="clearfix">
        &nbsp;
        <ReactTable
          data={datum.values}
          noDataText="No buckets available"
          columns={[
            {
              Header: 'Name',
              accessor: 'key',
              Cell: row => (<strong>{row.value}</strong>)
            },
            ...dates
          ]}
          defaultSorted={[{
            id: 'Cost',
            desc: true
          }]}
          defaultPageSize={10}
          className=" -highlight"
        />
      </div>
    );
  }

}

DifferentiatorChartComponent.propTypes = {
  values: PropTypes.object,
  interval: PropTypes.string.isRequired,
  legend: PropTypes.bool.isRequired,
  height: PropTypes.number.isRequired,
  margin: PropTypes.bool,
  table: PropTypes.bool
};

DifferentiatorChartComponent.defaultProps = {
  margin: true,
  table: false
};

export default DifferentiatorChartComponent;