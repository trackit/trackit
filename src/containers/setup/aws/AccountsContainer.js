import React, { Component } from 'react';
import { connect } from 'react-redux';

import Components from '../../../components';
import Actions from "../../../actions";
import PropTypes from "prop-types";
import s3square from '../../../assets/s3-square.png';

const List = Components.AWS.Accounts.List;
const Wizard = Components.AWS.Accounts.Wizard;
const Panel = Components.Misc.Panel;

// Accounts Container for AWS Accounts
export class AccountsContainer extends Component {

  componentWillMount() {
    this.props.getAccounts();
    this.props.newExternal();
  }

  render() {
    const noAccountsInfos = (
      <div>
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
      <Panel>

        <div>

          <h3 className="white-box-title no-padding inline-block">
            <img className="white-box-title-icon" src={s3square} alt="AWS square logo"/>
            AWS Accounts
          </h3>

          <div className="inline-block pull-right">
            <Wizard
              external={this.props.external}
              submitAccount={this.props.accountActions.new}
              clearAccount={this.props.accountActions.clearNew}
              submitBucket={this.props.newBill}
              account={this.props.newAccount}
            />
          </div>
          {
            (!this.props.accounts.length && this.props.match.params.hasAccounts === "false")
            && noAccountsInfos
          }

        </div>

        <List
          accounts={this.props.accounts}
          accountActions={this.props.accountActions}
        />

      </Panel>
    );
  }

}

AccountsContainer.propTypes = {
  accounts: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.number.isRequired,
      roleArn: PropTypes.string.isRequired,
      pretty: PropTypes.string
    })
  ),
  newAccount: PropTypes.shape({
    id: PropTypes.number.isRequired,
    roleArn: PropTypes.string.isRequired,
    pretty: PropTypes.string,
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
  newBill: PropTypes.func.isRequired,
  newExternal: PropTypes.func.isRequired
};

/* istanbul ignore next */
const mapStateToProps = (state) => ({
  accounts: state.aws.accounts.all,
  newAccount: state.aws.accounts.creation,
  external: state.aws.accounts.external
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getAccounts: () => {
    dispatch(Actions.AWS.Accounts.getAccounts())
  },
  accountActions: {
    new: (account) => {
      dispatch(Actions.AWS.Accounts.newAccount(account))
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
  newBill: (accountID, bill) => {
    dispatch(Actions.AWS.Accounts.newAccountBill(accountID, bill))
  },
  newExternal: () => {
    dispatch(Actions.AWS.Accounts.newExternal())
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(AccountsContainer);
