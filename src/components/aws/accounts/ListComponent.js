import React, { Component } from 'react';
import List, {
  ListItem,
  ListItemText,
} from 'material-ui/List';
import Spinner from 'react-spinkit';
import PropTypes from 'prop-types';
import Misc from '../../misc';
import Form from './FormComponent';
import Bills from './bills';

const DeleteConfirmation = Misc.DeleteConfirmation;

export class Item extends Component {

  constructor(props) {
    super(props);
    this.editAccount = this.editAccount.bind(this);
    this.deleteAccount = this.deleteAccount.bind(this);
  }

  editAccount = (body) => {
    this.setState({ editForm: false });
    console.log("Account edition is not available yet");
//    this.props.accountActions.edit(body);
  };

  deleteAccount = () => {
    console.log("Account deletion is not available yet");
//    this.props.accountActions.delete(this.props.account.id);
  };

  render() {
    return (
      <div>

        <ListItem divider>

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

    const accounts = (this.props.accounts.status && this.props.accounts.values.length ? (
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
        bills: PropTypes.arrayOf(
          PropTypes.shape({
            bucket: PropTypes.string.isRequired,
            path: PropTypes.string.isRequired
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
