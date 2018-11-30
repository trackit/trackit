import React, { Component } from 'react';
import { connect } from 'react-redux';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Spinner from 'react-spinkit';
import Misc from '../../../misc';
import PropTypes from 'prop-types';
import Form from './FormComponent';
import Actions from "../../../../actions";

const Dialog = Misc.Dialog;
const DeleteConfirmation = Misc.DeleteConfirmation;

const permissions = {
  2: "Read-Only",
  1: "Standard",
  0: "Administrator"
};

export class Item extends Component {

  constructor(props) {
    super(props);
    this.editAccountViewer = this.editAccountViewer.bind(this);
    this.deleteAccountViewer = this.deleteAccountViewer.bind(this);
  }

  editAccountViewer = (email, permission) => {
    this.props.editAccountViewer(this.props.account, this.props.accountViewer.sharedId, permission);
  };

  deleteAccountViewer = () => {
    this.props.deleteAccountViewer(this.props.account, this.props.accountViewer.sharedId);
  };

  render() {
    const permission = (<span className="badge blue-bg pull-right">{permissions[this.props.accountViewer.level]}</span>);

    return (
      <ListItem divider>

        <ListItemText
          disableTypography
          primary={
            <span>
              {this.props.accountViewer.email}
              {permission}
            </span>
          }
        />

        <div className="actions">

          <div className="inline-block">
            <Form
              account={this.props.account}
              AccountViewer={this.props.accountViewer}
              submit={this.editAccountViewer}
              status={this.props.editionStatus}
              clear={this.props.clearEdition}
              disabled={this.props.permissionLevel > this.props.accountViewer.level}
              permissionLevel={this.props.permissionLevel}
            />
          </div>
          &nbsp;
          <div className="inline-block">
            <DeleteConfirmation entity={`this shared account`} confirm={this.deleteAccountViewer} disabled={this.props.permissionLevel > this.props.accountViewer.level}/>
          </div>

        </div>

      </ListItem>
    );
  }

}

Item.propTypes = {
  account: PropTypes.number.isRequired,
  accountViewer: PropTypes.shape({
    sharedId: PropTypes.number.isRequired,
    email: PropTypes.string.isRequired,
    level: PropTypes.number.isRequired,
    userId: PropTypes.number.isRequired,
    sharingStatus: PropTypes.bool.isRequired
  }).isRequired,
  editionStatus: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    value: PropTypes.object
  }),
  editAccountViewer: PropTypes.func.isRequired,
  clearEdition: PropTypes.func.isRequired,
  deleteAccountViewer: PropTypes.func.isRequired,
  permissionLevel: PropTypes.number.isRequired
};

// List Component for AWS Accounts
export class ListComponent extends Component {

  constructor(props) {
    super(props);
    this.getAccountViewers = this.getAccountViewers.bind(this);
    this.newAccountViewer = this.newAccountViewer.bind(this);
    this.clearAccountViewers = this.clearAccountViewers.bind(this);
  }

  getAccountViewers() {
    this.props.getAccountViewers(this.props.account);
  }

  newAccountViewer(email, permission) {
    this.props.newAccountViewer(this.props.account, email, permission);
  }

  clearAccountViewers() {
    this.props.clearAccountViewers();
  }

  render() {
    const loading = (!this.props.accounts.status ? (<Spinner className="spinner" name='circle'/>) : null);

    const error = (this.props.accounts.error ? ` (${this.props.accounts.error.message})` : null);
    const noAccounts = (this.props.accounts.status && (!this.props.accounts.values || !this.props.accounts.values.length || error) ? <div className="alert alert-warning" role="alert">This account is not shared with other users.{error}</div> : "");

    const accounts = (this.props.accounts.status && this.props.accounts.values && this.props.accounts.values.length ? (
      this.props.accounts.values.map((account, index) => (
        <Item
          key={index}
          accountViewer={account}
          account={this.props.account}
          editAccountViewer={this.props.editAccountViewer}
          editionStatus={this.props.accountViewerEdition}
          clearEdition={this.props.clearEditAccountViewer}
          deleteAccountViewer={this.props.deleteAccountViewer}
          permissionLevel={this.props.permissionLevel}
         />
      ))
    ) : null);

    const form = (<Form
      account={this.props.account}
      submit={this.newAccountViewer}
      status={this.props.accountViewerCreation}
      clear={this.props.clearNewAccountViewer}
      permissionLevel={this.props.permissionLevel}
    />);

    return (
      <Dialog
        buttonName={<span><i className="fa fa-share-alt"/> Share</span>}
        disabled={this.props.disabled}
        title="Share this account"
        secondActionName="Close"
        onOpen={this.getAccountViewers}
        onClose={this.clearAccountViewers}
        titleChildren={form}
      >

        <List className="bills-list">
          {loading}
          {noAccounts}
          {accounts}
        </List>

      </Dialog>
    );
  }

}

ListComponent.propTypes = {
  account: PropTypes.number.isRequired,
  accounts: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    values: PropTypes.arrayOf(
      PropTypes.shape({
        sharedId: PropTypes.number.isRequired,
        email: PropTypes.string.isRequired,
        level: PropTypes.number.isRequired,
        userId: PropTypes.number.isRequired,
        sharingStatus: PropTypes.bool.isRequired
      }).isRequired,
    ),
  }),
  accountViewerCreation: PropTypes.object,
  accountViewerEdition: PropTypes.object,
  clearAccountViewers: PropTypes.func.isRequired,
  getAccountViewers: PropTypes.func.isRequired,
  newAccountViewer: PropTypes.func.isRequired,
  editAccountViewer: PropTypes.func.isRequired,
  clearNewAccountViewer: PropTypes.func.isRequired,
  clearEditAccountViewer: PropTypes.func.isRequired,
  deleteAccountViewer: PropTypes.func.isRequired,
  permissionLevel: PropTypes.number.isRequired,
  disabled: PropTypes.bool,
};

ListComponent.defaultProps = {
  disabled: false
};

/* istanbul ignore next */
const mapStateToProps = (state) => ({
  accounts: state.aws.accounts.accountViewers,
  accountViewerCreation: state.aws.accounts.addAccountViewer,
  accountViewerEdition: state.aws.accounts.editAccountViewer,
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getAccountViewers: (accountID) => {
    dispatch(Actions.AWS.Accounts.getAccountViewers(accountID));
  },
  clearAccountViewers: () => {
    dispatch(Actions.AWS.Accounts.clearAccountViewers());
  },
  newAccountViewer: (accountID, email, permission) => {
    dispatch(Actions.AWS.Accounts.newAccountViewer(accountID, email, permission));
  },
  editAccountViewer: (accountID, shareID, permission) => {
    dispatch(Actions.AWS.Accounts.editAccountViewer(accountID, shareID, permission));
  },
  clearNewAccountViewer: () => {
    dispatch(Actions.AWS.Accounts.clearNewAccountViewer());
  },
  clearEditAccountViewer: () => {
    dispatch(Actions.AWS.Accounts.clearEditAccountViewer());
  },
  deleteAccountViewer: (accountID, shareID) => {
    dispatch(Actions.AWS.Accounts.deleteAccountViewer(accountID, shareID));
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(ListComponent);
