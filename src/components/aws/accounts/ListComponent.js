import React, { Component } from 'react';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Spinner from 'react-spinkit';
import PropTypes from 'prop-types';
import Misc from '../../misc';
import Form from './FormComponent';
import Bills from './bills';
import TeamSharing from './teamSharing';
import Status from '../../../common/awsAccountStatus';

const Tooltip = Misc.Popover;
const DeleteConfirmation = Misc.DeleteConfirmation;

export class Item extends Component {

  constructor(props) {
    super(props);
    this.editAccount = this.editAccount.bind(this);
    this.deleteAccount = this.deleteAccount.bind(this);
  }

  editAccount = (body) => {
    this.setState({ editForm: false });
    body.id = this.props.account.id;
    this.props.accountActions.edit(body);
  };

  deleteAccount = () => {
    this.props.accountActions.delete(this.props.account.id);
  };

  render() {
    const status = Status.getAWSAccountStatus(this.props.account);
    const accountBadge = Status.getBadge(status);
    const infoBanner = Status.getInformationBanner(status);
    const formattedInfoBanner = (infoBanner ? (
      <ListItem className="account-alert">
        {infoBanner}
      </ListItem>
    ) : null);

    const subAccounts = (this.props.account.subAccounts && this.props.account.subAccounts.length ? (
      <ListItem>
        <List disablePadding className="account-subaccounts">
          {this.props.account.subAccounts.map((subAccount, key) => (
            <Item
              key={key}
              account={subAccount}
              accountActions={this.props.accountActions}
              subAccount
            />
          ))}
        </List>
      </ListItem>
    ) : null);

    return (
      <div>

        <ListItem className="account-item">

          {accountBadge}

          <ListItemText
            disableTypography
            className="account-name"
            primary={this.props.account.pretty || this.props.account.roleArn}
          />

          <div className="actions">

            <div className="inline-block">
              {(!this.props.account.accountOwner ? (
                <Tooltip icon={(
                  <div className="btn btn=default">
                    <i className="fa account-badge fa-share-alt"/>
                  </div>
                )} tooltip="This account is shared by another user" placement="left"/>
              ) : null)}
            </div>

            <div className="inline-block">
              <TeamSharing.List account={this.props.account.id} disabled={this.props.account.permissionLevel === 2} permissionLevel={this.props.account.permissionLevel}/>
            </div>
            &nbsp;
            {!this.props.subAccount ? (
              <div>
                <div className="inline-block">
                  <Bills.List account={this.props.account.id} disabled={this.props.account.permissionLevel !== 0}/>
                </div>
                &nbsp;
              </div>
            ) : null}
            <div className="inline-block">
              <Form
                account={this.props.account}
                submit={this.editAccount}
                disabled={this.props.account.permissionLevel !== 0}
                subAccount={this.props.subAccount}
              />
            </div>
            &nbsp;
            <div className="inline-block">
              <DeleteConfirmation entity="account" confirm={this.deleteAccount} disabled={this.props.account.permissionLevel !== 0}/>
            </div>

          </div>

        </ListItem>

        {formattedInfoBanner}
        {subAccounts}

      </div>
    );
  }

}

Item.propTypes = {
  account: PropTypes.shape({
    id: PropTypes.number.isRequired,
    roleArn: PropTypes.string.isRequired,
    pretty: PropTypes.string,
    permissionLevel: PropTypes.number,
    payer: PropTypes.bool.isRequired,
    status: PropTypes.shape({
      value: PropTypes.string.isRequired,
      detail: PropTypes.string.isRequired,
    }),
    billRepositories: PropTypes.arrayOf(
      PropTypes.shape({
        error: PropTypes.string.isRequired,
        accountOwner: PropTypes.bool,
        nextPending: PropTypes.bool.isRequired,
        bucket: PropTypes.string.isRequired,
        prefix: PropTypes.string.isRequired,
        status: PropTypes.shape({
          value: PropTypes.string.isRequired,
          detail: PropTypes.string.isRequired,
        })
      })
    ),
  }),
  accountActions: PropTypes.shape({
    edit: PropTypes.func.isRequired,
    delete: PropTypes.func.isRequired,
  }).isRequired,
  subAccount: PropTypes.bool
};

Item.defaultProps = {
  subAccount: false
};

// List Component for AWS Accounts
class ListComponent extends Component {

  render() {
    const loading = (!this.props.accounts.status ? (<Spinner className="spinner" name='circle'/>) : null);

    const error = (this.props.accounts.error ? ` (${this.props.accounts.error.message})` : null);
    const noAccounts = (this.props.accounts.status && (!this.props.accounts.values || !this.props.accounts.values.length || error) ? <div className="alert alert-warning" role="alert">No account available{error}</div> : "");

    const accounts = (this.props.accounts.status && this.props.accounts.values && this.props.accounts.values.length ? (
      this.props.accounts.values.map((account, index) => (
        <div className="white-box" key={index}>
          <Item
            account={account}
            accountActions={this.props.accountActions}
          />
        </div>
      ))
    ) : null);

    return (
      <List disablePadding className="accounts-list">
        {loading}
        {noAccounts}
        {accounts}
      </List>
    );
  }

}

ListComponent.propTypes = {
  accounts: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    values: PropTypes.arrayOf(
      PropTypes.shape({
        id: PropTypes.number.isRequired,
        roleArn: PropTypes.string.isRequired,
        pretty: PropTypes.string,
        permissionLevel: PropTypes.number,
        payer: PropTypes.bool.isRequired,
        billRepositories: PropTypes.arrayOf(
          PropTypes.shape({
            error: PropTypes.string.isRequired,
            nextPending: PropTypes.bool.isRequired,
            bucket: PropTypes.string.isRequired,
            prefix: PropTypes.string.isRequired
          })
        ),
      })
    ),
  }),
  accountActions: PropTypes.shape({
    edit: PropTypes.func.isRequired,
    delete: PropTypes.func.isRequired,
  }).isRequired,
};

export default ListComponent;
