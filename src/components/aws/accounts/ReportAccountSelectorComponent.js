import React, {Component} from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

import Actions from "../../../actions";
import Selector from '../../misc/Selector';

export class ReportAccountSelectorComponent extends Component {

  render() {
    const listedAccounts = (this.props.accounts.values && this.props.accounts.values.length > 0 ? this.props.accounts.values : null);
    const availableAccounts = {};
    if (!listedAccounts) {
      return null;
    }
    listedAccounts.forEach((account) => {
      availableAccounts[account.id] = account.pretty;
    });
    return(
      <Selector values={availableAccounts} selected={this.props.account} selectValue={this.props.selectAccount}/>
    );
  }

}

ReportAccountSelectorComponent.propTypes = {
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
  selectAccount: PropTypes.func.isRequired,
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  accounts: aws.accounts.all,
  account: aws.reports.account,
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  selectAccount: (accountId) => {
    dispatch(Actions.AWS.Reports.selectAccount(accountId));
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(ReportAccountSelectorComponent);
