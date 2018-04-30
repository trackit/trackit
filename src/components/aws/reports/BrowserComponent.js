import React, {Component} from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import ReactTable from 'react-table';
import Spinner from 'react-spinkit';

import Actions from "../../../actions";


export class BrowserComponent extends Component {
  render() {
    if (!this.props.reportList.status) {
      return (<Spinner className="spinner" name='circle'/>);
    }
    const error = (this.props.reportList.error ? ` (${this.props.reportList.error.message})` : null);
    const noReports = (this.props.reportList.status && (!this.props.reportList.values || !this.props.reportList.values.length || error) ? <div className="alert alert-warning" role="alert">No reports available{error}</div> : "");
    if (noReports !== '') {
      return (noReports);
    }
    const data = (this.props.reportList.status && this.props.reportList.values && this.props.reportList.values.length ? (
      this.props.reportList.values.map((report, index) => (
        {Name: report}
      ))
    ) : []);
    return(
      <ReactTable
        data={data}
        noDataText="No reports available"
        columns={[
            {
              Header: 'Reports',
              accessor: 'Name',
            },
          ]
        }
        defaultPageSize={10}
        defaultSorted={[{
          id: 'Name',
          desc: true
        }]}
        className=" -highlight"
        getTdProps={(state, rowInfo, column, instance) => {
          return {
            onClick: (e, handleOriginal) => {
              let res = rowInfo.original.Name.split('/');
              this.props.startDownload(this.props.account, res[0], res[1]);
              /* istanbul ignore next */
              if (handleOriginal) {
                handleOriginal();
              }
            }
          };
        }}
      />
    );
  }
}

BrowserComponent.propTypes = {
  reportList: PropTypes.object.isRequired,
  startDownload: PropTypes.func.isRequired,
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  account: aws.reports.account,
  reportList: aws.reports.reportList
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  startDownload: (accountId, reportType, fileName) => {
    dispatch(Actions.AWS.Reports.requestDownloadReport(accountId, reportType, fileName));
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(BrowserComponent);
