import React, { Component } from 'react';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Spinner from 'react-spinkit';
import PropTypes from 'prop-types';
import Misc from '../../misc';
import Form from './FormComponent';
import Bills from './bills';

const DeleteConfirmation = Misc.DeleteConfirmation;
const Popover = Misc.Popover;

export class Item extends Component {

  constructor(props) {
    super(props);
    this.editAccount = this.editAccount.bind(this);
    this.deleteAccount = this.deleteAccount.bind(this);
    this.getAccountBadge = this.getAccountBadge.bind(this);
  }

  editAccount = (body) => {
    this.setState({ editForm: false });
    body.id = this.props.account.id;
    this.props.accountActions.edit(body);
  };

  deleteAccount = () => {
    this.props.accountActions.delete(this.props.account.id);
  };

  getAccountBadge = () => {
    let error = false;
    let pending = false;
    this.props.account.billRepositories.forEach((billRepository) => {
      if (billRepository.error !== "")
        error = true;
      if (billRepository.nextPending)
        pending = true;
    });
    if (error || !this.props.account.billRepositories.length)
      return (
          <Popover
            children={<i className="fa account-badge fa-times-circle"/>}
            popOver={"Please check your bill locations"}
          />
      );
    else if (pending)
      return (
          <Popover
            children={<i className="fa account-badge fa-clock-o"/>}
            popOver={"Import in progress"}
          />
      );
    return (<i className="fa account-badge fa-check-circle"/>);
  };

  render() {
    return (
      <div>

        <ListItem divider>

          {this.getAccountBadge()}
          <ListItemText
            disableTypography
            className="account-name"
            primary={this.props.account.pretty || this.props.account.roleArn}
          />

          <div className="actions">

            <div className="inline-block">
              <Bills.List account={this.props.account.id} />
            </div>
            &nbsp;
            <div className="inline-block">
              <Form
                account={this.props.account}
                submit={this.editAccount}
              />
            </div>
            &nbsp;
            <div className="inline-block">
              <DeleteConfirmation entity="account" confirm={this.deleteAccount}/>
            </div>

          </div>

        </ListItem>

      </div>
    );
  }

}

Item.propTypes = {
  account: PropTypes.shape({
    id: PropTypes.number.isRequired,
    roleArn: PropTypes.string.isRequired,
    pretty: PropTypes.string,
    billRepositories: PropTypes.arrayOf(
      PropTypes.shape({
        error: PropTypes.string.isRequired,
        nextPending: PropTypes.bool.isRequired,
        bucket: PropTypes.string.isRequired,
        prefix: PropTypes.string.isRequired
      })
    ),
  }),
  accountActions: PropTypes.shape({
    edit: PropTypes.func.isRequired,
    delete: PropTypes.func.isRequired,
  }).isRequired,
};

// List Component for AWS Accounts
class ListComponent extends Component {

  render() {
    const loading = (!this.props.accounts.status ? (<Spinner className="spinner" name='circle'/>) : null);

    const error = (this.props.accounts.error ? ` (${this.props.accounts.error.message})` : null);
    const noAccounts = (this.props.accounts.status && (!this.props.accounts.values || !this.props.accounts.values.length || error) ? <div className="alert alert-warning" role="alert">No account available{error}</div> : "");

    const accounts = (this.props.accounts.status && this.props.accounts.values && this.props.accounts.values.length ? (
      this.props.accounts.values.map((account, index) => (
        <Item
          key={index}
          account={account}
          accountActions={this.props.accountActions}
        />
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
