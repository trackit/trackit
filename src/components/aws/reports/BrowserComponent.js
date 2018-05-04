import React, {Component} from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import Spinner from 'react-spinkit';
import List, { ListItem, ListItemText } from 'material-ui/List';

import Actions from "../../../actions";


export class BrowserComponent extends Component {

  constructor(props) {
    super(props);
    this.handleDownloadClick = this.handleDownloadClick.bind(this);
  }

  handleDownloadClick = (e, report) => {
    e.preventDefault();
    let reportInfos = report.split('/');
    this.props.startDownload(this.props.account, reportInfos[0], reportInfos[1]);
  };

  render() {
    if (!this.props.reportList.status) {
      return (<Spinner className="spinner" name='circle'/>);
    }
    const error = (this.props.reportList.error ? ` (${this.props.reportList.error.message})` : null);
    const noReports = (this.props.reportList.status && (!this.props.reportList.values || !this.props.reportList.values.length || error) ? <div className="alert alert-warning" role="alert">No reports available{error}</div> : "");
    if (noReports !== '') {
      return (noReports);
    }
    const listItems = this.props.reportList.values.map((report, index) => (
      <div key={report}>
        <ListItem divider>
          <i class="fa fa-file-excel-o fa-2x red-color"></i>
          <ListItemText
            disableTypography
            className="report-name"
            primary={report}
          />
          <div className="actions">
            <div className="inline-block">
              <button className="btn btn-default" onClick={(e) => this.handleDownloadClick(e, report)}>Download</button>
            </div>
          </div>
        </ListItem>
      </div>
    ));
    return(
      <List disablePadding className="reports-list">
        {listItems}
      </List>
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
