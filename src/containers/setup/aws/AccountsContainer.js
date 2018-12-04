import React, { Component } from 'react';
import { connect } from 'react-redux';

import Components from '../../../components';
import Actions from "../../../actions";
import PropTypes from "prop-types";

const List = Components.AWS.Accounts.List;
const Wizard = Components.AWS.Accounts.Wizard;
const Status = Components.AWS.Accounts.Bills.Status;

// Accounts Container for AWS Accounts
export class AccountsContainer extends Component {

  componentWillMount() {
    this.props.getAccounts();
    this.props.newExternal();
  }

  render() {
    const noAccountsInfos = (
      <div id="welcome">
        <hr />
        <div className="alert alert-info" role="alert" style={{ fontSize: '15px', lineHeight: '2' }}>
          <strong>
            <i className="fa fa-info-circle"/>
            &nbsp;
            Welcome to TrackIt !
          </strong>
          <br />
          {"It seems you don't have any AWS account setup yet."}
          <br />
          {"Please click the above "}
          <strong>Add</strong>
          {" button and follow the instructions. It will get you up and running in no time !"}
          <br />
          {"Thank you for using TrackIt !"}
        </div>
      </div>
    );

    return (
      <div>

        <div className="white-box">

          <h3 className="white-box-title no-padding inline-block">
            {/* <img className="white-box-title-icon" src={s3square} alt="AWS square logo"/> */}
            <i className="fa fa-amazon"></i>
            &nbsp;
            AWS Accounts
          </h3>

          <div className="inline-block pull-right">
            <div className="inline-block">
              {
                !(!this.props.accounts.length && this.props.match.params.hasAccounts === "false") && (
                  <Status
                    bills={this.props.billsStatus}
                    billsStatusActions={this.props.billsStatusActions}
                  />
                )
              }
            </div>
            &nbsp;
            <div className="inline-block">
              <Wizard
                external={this.props.external}
                submitAccount={this.props.accountActions.new}
                clearAccount={this.props.accountActions.clearNew}
                submitBucket={this.props.addBill}
                clearBucket={this.props.clearBill}
                account={this.props.newAccount}
                bill={this.props.newBill}
              />
            </div>
          </div>

          {
            (!(this.props.accounts.values && this.props.accounts.values.length) && this.props.match.params.hasAccounts === "false")
            && noAccountsInfos
          }

        </div>

        <List
          accounts={this.props.accounts}
          accountActions={this.props.accountActions}
        />

      </div>
    );
  }

}

AccountsContainer.propTypes = {
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
  newAccount: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    value: PropTypes.shape({
      id: PropTypes.number.isRequired,
      roleArn: PropTypes.string.isRequired,
      pretty: PropTypes.string
    })
  }),
  newBill: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error)
  }),
  external: PropTypes.shape({
    external: PropTypes.string.isRequired,
    accountId: PropTypes.string.isRequired,
  }),
  getAccounts: PropTypes.func.isRequired,
  accountActions: PropTypes.shape({
    new: PropTypes.func.isRequired,
    edit: PropTypes.func.isRequired,
    delete: PropTypes.func.isRequired,
  }).isRequired,
  billsStatus: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    values: PropTypes.arrayOf(
      PropTypes.shape({
        BillRepositoryId: PropTypes.number.isRequired,
        AwsAccountPretty: PropTypes.string.isRequired,
        AwsAccountId: PropTypes.number.isRequired,
        bucket: PropTypes.string.isRequired,
        prefix: PropTypes.string.isRequired,
        nextStarted: PropTypes.string.isRequired,
        nextPending: PropTypes.bool.isRequired,
        lastStarted: PropTypes.string.isRequired,
        lastFinished: PropTypes.string.isRequired,
        lastError: PropTypes.string.isRequired
      })
    )
  }),
  billsStatusActions: PropTypes.shape({
    get: PropTypes.func.isRequired,
    clear: PropTypes.func.isRequired,
  }).isRequired,
  addBill: PropTypes.func.isRequired,
  clearBill: PropTypes.func.isRequired,
  newExternal: PropTypes.func.isRequired
};

/* istanbul ignore next */
const mapStateToProps = (state) => ({
  accounts: state.aws.accounts.all,
  newAccount: state.aws.accounts.creation,
  newBill: state.aws.accounts.billCreation,
  external: state.aws.accounts.external,
  billsStatus: state.aws.accounts.billsStatus
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getAccounts: () => {
    dispatch(Actions.AWS.Accounts.getAccounts())
  },
  accountActions: {
    new: (account, bill) => {
      dispatch(Actions.AWS.Accounts.newAccount(account, bill))
    },
    clearNew: () => {
      dispatch(Actions.AWS.Accounts.clearNewAccount());
    },
    edit: (account) => {
      dispatch(Actions.AWS.Accounts.editAccount(account))
    },
    delete: (accountID) => {
      dispatch(Actions.AWS.Accounts.deleteAccount(accountID));
    },
  },
  billsStatusActions: {
    get: () => {
      dispatch(Actions.AWS.Accounts.getAccountBillsStatus());
    },
    clear: () => {
      dispatch(Actions.AWS.Accounts.clearAccountBillsStatus());
    }
  },
  addBill: (accountID, bill) => {
    dispatch(Actions.AWS.Accounts.newAccountBill(accountID, bill))
  },
  clearBill: () => {
    dispatch(Actions.AWS.Accounts.clearNewAccountBill())
  },
  newExternal: () => {
    dispatch(Actions.AWS.Accounts.newExternal())
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(AccountsContainer);
