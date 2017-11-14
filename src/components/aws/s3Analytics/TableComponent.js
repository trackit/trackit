import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ReactTable from "react-table";

import { formatBytes, formatPrice } from '../../../common/formatters';


// S3AnalyticsTableComponent Component
class TableComponent extends Component {

    render() {
      return (
        <div>
          <ReactTable
            data={this.props.data}
            columns={
              [
                {
                  Header: 'Name',
                  accessor: '_id',
                  Cell: row => (
                    <strong>{row.value}</strong>
                  )
                },
                {
                  Header: 'Size',
                  accessor: 'size',
                  Cell: row => (
                    formatBytes(row.value, 1)
                  )
                },
                {
                  Header: 'Cost',
                  columns: [
                    {
                      Header: 'Storage',
                      accessor: 'storage_cost',
                      Cell: row => (
                        formatPrice(row.value)
                      )
                    },
                    {
                      Header: 'Bandwidth',
                      accessor: 'bw_cost',
                      Cell: row => (
                        formatPrice(row.value)
                      )
                    },
                    {
                      Header: 'Total',
                      accessor: 'total_cost',
                      Cell: row => (
                        <span className="total-cell">{formatPrice(row.value)}</span>
                      )
                    },
                  ]
                },
                {
                  Header: 'Data transfers',
                  columns: [
                    {
                      Header: 'In',
                      accessor: 'transfer_in',
                      Cell: row => (
                        formatBytes(row.value)
                      )
                    },
                    {
                      Header: 'Out',
                      accessor: 'transfer_out',
                      Cell: row => (
                        formatBytes(row.value)
                      )
                    },
                  ]
                },
                {
                  Header: 'Chargify',
                  accessor: 'chargify',
                  Cell: row => (
                    <span>
                      <span style={{
                        color: row.value === 'not_synced' ? '#ff2e00'
                          : row.value === 'in_sync' ? '#ffbf00'
                          : '#57d500',
                        transition: 'all .3s ease'
                      }}>
                        &#x25cf;
                      </span>
                      &nbsp;
                      {row.value}
                    </span>
                  )
                },


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
  data: PropTypes.arrayOf(
    PropTypes.shape({
      _id: PropTypes.string.isRequired,
      size: PropTypes.number.isRequired,
      storage_cost: PropTypes.number.isRequired,
      bw_cost: PropTypes.number.isRequired,
      total_cost: PropTypes.number.isRequired,
      transfer_in: PropTypes.number.isRequired,
      transfer_out: PropTypes.number.isRequired,
      chargify: PropTypes.oneOf(['not_synced', 'in_sync', 'synced'])
    })
  ),
};

export default TableComponent;
