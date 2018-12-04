import React, { Component } from 'react';
import {Redirect} from "react-router-dom";
import { connect } from 'react-redux';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogActions from '@material-ui/core/DialogActions';
import List from '@material-ui/core/List';
import PropTypes from 'prop-types';
import Actions from "../../../../actions/index";
import Spinner from 'react-spinkit';
import Item from './ItemComponent';

const styles = {
  badge : {
    fontSize: '14px',
    fontWeight: '500',
    cursor: "pointer"
  },
  icon: {
    fontSize: '16px',
  }
};

const schedulerDuration = 1000 * 60 * 30; // 30 minutes

// Accounts Status Component for AWS Account status
class AccountsStatusComponent extends Component {

  constructor(props) {
    super(props);
    this.state = {
      open: false,
      goToSetup: false,
    };
    this.openDialog = this.openDialog.bind(this);
    this.closeDialog = this.closeDialog.bind(this);
    this.setupAccounts = this.setupAccounts.bind(this);
    this.scheduler = null;
  }

  componentDidMount() {
    this.props.getAccounts();
    this.setScheduler();
  }

  componentWillReceiveProps(nextProps) {
    clearTimeout(this.scheduler);
    if (this.props.accounts !== nextProps.accounts && nextProps.accounts.status && nextProps.accounts.hasOwnProperty("values"))
      this.setScheduler();
  }

  setScheduler() {
    this.scheduler = setTimeout(() => {
      this.props.getAccounts();
    }, schedulerDuration);
  }

  openDialog = (e) => {
    e.preventDefault();
    this.setState({open: true});
  };

  closeDialog = (e) => {
    e.preventDefault();
    this.setState({
      open: false,
    });
  };

  setupAccounts = (e) => {
    e.preventDefault();
    this.setState({
      open: false,
      goToSetup: true,
    });
  };

  getText = () => {
    const error = (this.props.accounts.error ? ` (${this.props.accounts.error.message})` : '');

    if (!this.props.accounts.status)
      return null;

    if (!this.props.accounts.values || !this.props.accounts.values.length || error)
      return `No accounts ${error}`;

    let accountsNumber = this.props.accounts.values.length;
    this.props.accounts.values.forEach((account) => {
      if (account.hasOwnProperty("subAccounts") && account.subAccounts)
        accountsNumber += account.subAccounts.length;
    });

    if (this.props.selection.length === 0 || accountsNumber === this.props.selection.length)
      return `Displaying All accounts`;
    if (this.props.selection.length === 1)
      return `Displaying ${this.props.selection[0].pretty}`;
    return `Displaying ${this.props.selection.length} accounts`;
  };

  render() {
    if (this.state.goToSetup) {
      this.setState({goToSetup: false});
      return (<Redirect to="/app/setup"/>);
    }

    const isSelected = (item) => (this.props.selection.find((value) => (value.id === item.id)) !== undefined);

    const loading = (!this.props.accounts.status ? (<Spinner className="spinner" name='circle'/>) : null);

    const error = (this.props.accounts.error ? ` (${this.props.accounts.error.message})` : null);
    const noAccounts = (this.props.accounts.status && (!this.props.accounts.values || !this.props.accounts.values.length || error) ? <div className="alert alert-warning" role="alert">No account available{error}</div> : "");

    const accounts = (this.props.accounts.status && this.props.accounts.values && this.props.accounts.values.length ? (
      this.props.accounts.values.map((account, index) => (
        <Item
          key={index}
          account={account}
          select={this.props.select}
          isSelected={isSelected}
        />
      ))
    ) : null);

    return (
      <div>

        <span className="badge" onClick={this.openDialog} style={styles.badge}>
          <span><i className="fa fa-amazon" style={styles.icon}/>&nbsp;&nbsp;</span>
          {this.getText()}
        </span>

        <Dialog open={this.state.open} fullWidth>

          <DialogTitle disableTypography>
            <h1>AWS Accounts</h1>
          </DialogTitle>

          <DialogContent>

            <List disablePadding className="accounts-list">
              {loading}
              {noAccounts}
              {accounts}
            </List>

            <DialogActions>

              <button className="btn btn-default pull-left" onClick={this.setupAccounts}>
                Setup Accounts
              </button>


              <button className="btn btn-default pull-left" onClick={this.closeDialog}>
                Close
              </button>

            </DialogActions>

          </DialogContent>

        </Dialog>
      </div>
    );
  }

}

AccountsStatusComponent.propTypes = {
  accounts: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    values: PropTypes.arrayOf(
      PropTypes.shape({
        id: PropTypes.number.isRequired,
        accountOwner: PropTypes.bool.isRequired,
        awsIdentity: PropTypes.string.isRequired,
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
        subAccounts: PropTypes.arrayOf(
          PropTypes.shape({
            id: PropTypes.number.isRequired,
            accountOwner: PropTypes.bool.isRequired,
            awsIdentity: PropTypes.string.isRequired,
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
        )
      })
    ),
  }),
  selection: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.number.isRequired,
      awsIdentity: PropTypes.string.isRequired,
      pretty: PropTypes.string,
    })
  ),
  select: PropTypes.func.isRequired,
  clear: PropTypes.func.isRequired,
  getAccounts: PropTypes.func.isRequired
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  accounts: aws.accounts.all,
  selection: aws.accounts.selection,
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getAccounts: () => {
    dispatch(Actions.AWS.Accounts.getAccounts());
  },
  select: (account) => {
    dispatch(Actions.AWS.Accounts.selectAccount(account));
  },
  clear: () => {
    dispatch(Actions.AWS.Accounts.clearAccounts());
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(AccountsStatusComponent);
