import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ReactTable from "react-table";
import Spinner from 'react-spinkit';

import {formatGigaBytes, formatPrice} from '../../../common/formatters';


// S3AnalyticsTableComponent Component
class TableComponent extends Component {

  render() {
    if (!this.props.data || !this.props.data.status)
      return (<Spinner className="spinner clearfix" name='circle'/>);

    if (this.props.data && this.props.data.status && this.props.data.hasOwnProperty("error"))
      return (<div className="alert alert-warning" role="alert">Data not available ({this.props.data.error.message})</div>);

    const data = Object.keys(this.props.data.values).map((id) => ({
      id,
      ...this.props.data.values[id],
      TotalCost: (this.props.data.values[id].StorageCost + this.props.data.values[id].BandwidthCost)
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
                    Cell: row => (formatGigaBytes(row.value))
                  }, {
                    Header: 'Out',
                    accessor: 'DataOut',
                    Cell: row => (formatGigaBytes(row.value))
                  }
                ]
              }
            ]
          }
          defaultSorted={[{
            id: 'TotalCost',
            desc: true
          }]}
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
