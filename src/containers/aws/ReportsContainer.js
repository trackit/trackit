import React, {Component} from 'react';
import {connect} from 'react-redux';
import PropTypes from 'prop-types';

import Actions from "../../actions";
import Components from '../../components';
import s3square from '../../assets/s3-square.png';

const Panel = Components.Misc.Panel;
const SingleAccountSelector = Components.AWS.Accounts.SingleAccountSelector;
const ReportsBrowser = Components.AWS.Reports.Browser;

// S3AnalyticsContainer Component
export class ReportsContainer extends Component {

  componentWillMount() {
    if (!this.props.accounts.status) {
      this.props.getAccounts();
    } else {
      this.props.requestGetReports(this.props.account);
    }
  }

  componentDidUpdate() {
    if (this.props.account !== '' && !this.props.reportList.status) {
      this.props.requestGetReports(this.props.account);
    }
  }

  render() {
    const error = (this.props.downloadStatus.error ? ` (${this.props.downloadStatus.error.message})` : null);
    const downloadError = (this.props.downloadStatus.failed ? <div className="no-padding"><div className="alert alert-warning" role="alert">Failed to download report{error}</div></div> : null);
    return (
      <Panel>
          {downloadError}
          <div className="clearfix">
            <h3 className="white-box-title no-padding inline-block">
              <img className="white-box-title-icon" src={s3square} alt="AWS square logo"/>
              AWS Reports
            </h3>
            <div className="inline-block pull-right">
              Selected account :
              <SingleAccountSelector/>
            </div>
          </div>

          <div className="no-padding">
            <ReportsBrowser/>
          </div>
      </Panel>
    );
  }
}

ReportsContainer.propTypes = {
  accounts: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    values: PropTypes.arrayOf(
      PropTypes.shape({
        id: PropTypes.number.isRequired,
        roleArn: PropTypes.string.isRequired,
        pretty: PropTypes.string,
        bills: PropTypes.arrayOf(
          PropTypes.shape({
            bucket: PropTypes.string.isRequired,
            path: PropTypes.string.isRequired
          })
        ),
      })
    ),
  }),
  account: PropTypes.string,
  reportList: PropTypes.object.isRequired,
  downloadStatus: PropTypes.object.isRequired,
  getAccounts: PropTypes.func.isRequired,
  selectAccount: PropTypes.func.isRequired,
  requestGetReports: PropTypes.func.isRequired,
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  accounts: aws.accounts.all,
  account: aws.reports.account,
  reportList: aws.reports.reportList,
  downloadStatus: aws.reports.download
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getAccounts: () => {
    dispatch(Actions.AWS.Accounts.getAccounts());
  },
  selectAccount: (accountId) => {
    dispatch(Actions.AWS.Reports.selectAccount(accountId));
  },
  requestGetReports: (accountId) => {
    dispatch(Actions.AWS.Reports.requestGetReports(accountId));
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(ReportsContainer);
