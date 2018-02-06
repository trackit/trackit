import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ReactTable from "react-table";

import {formatBytes, formatGigaBytes, formatPrice} from '../../../common/formatters';


// S3AnalyticsTableComponent Component
class TableComponent extends Component {

  render() {
    const data = Object.keys(this.props.data).map((id) => ({
      id,
      ...this.props.data[id],
      TotalCost: (this.props.data[id].StorageCost + this.props.data[id].BandwidthCost)
    }));

    /* istanbul ignore next */
    return (
      <div>
        <ReactTable
          data={data}
          noDataText="No buckets available"
          columns={[
              {
                Header: 'Name',
                accessor: 'id',
                Cell: row => (<strong>{row.value}</strong>)
              }, {
                Header: 'Size',
                accessor: 'GbMonth',
                Cell: row => (formatGigaBytes(row.value, 1))
              }, {
                Header: 'Cost',
                columns: [
                  {
                    Header: 'Storage',
                    accessor: 'StorageCost',
                    Cell: row => (formatPrice(row.value))
                  }, {
                    Header: 'Bandwidth',
                    accessor: 'BandwidthCost',
                    Cell: row => (formatPrice(row.value))
                  }, {
                    Header: 'Total',
                    accessor: 'TotalCost',
                    Cell: row => (<span className="total-cell">{formatPrice(row.value)}</span>)
                  }
                ]
              }, {
                Header: 'Data transfers',
                columns: [
                  {
                    Header: 'In',
                    accessor: 'DataIn',
                    Cell: row => (formatBytes(row.value))
                  }, {
                    Header: 'Out',
                    accessor: 'DataOut',
                    Cell: row => (formatBytes(row.value))
                  }
                ]
              }
            ]
          }
          defaultPageSize={10}
          className=" -highlight"
        />
      </div>
    );
  }

}

TableComponent.propTypes = {
  data: PropTypes.object
};

export default TableComponent;
