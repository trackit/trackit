import React, { Component } from 'react';
import { connect } from 'react-redux';

import Components from '../../../components';
import Actions from "../../../actions";
import PropTypes from "prop-types";

const List = Components.AWS.Accounts.List;
const Form = Components.AWS.Accounts.Form;

// Accounts Container for AWS Accounts
class AccountsContainer extends Component {

  constructor(props) {
    super(props);
    this.addAccount = this.addAccount.bind(this);
  }

  componentWillMount() {
    this.props.getAccounts();
    this.props.newExternal();
  }

  addAccount = (account) => {
    this.props.newAccount(account);
    this.props.newExternal();
  };

  render() {
    return (
      <div className="panel panel-default">
        <div className="panel-heading">
          <h3 className="panel-title">AWS Accounts</h3>
        </div>
        <div className="panel-body">
          <List accounts={this.props.accounts}/>
          <Form submit={this.addAccount} external={this.props.external}/>
        </div>
      </div>
    );
  }

}

AccountsContainer.propTypes = {
  accounts: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.number.isRequired,
      roleArn: PropTypes.string.isRequired,
      userId: PropTypes.number.isRequired,
      pretty: PropTypes.string.isRequired
    })
  ),
  external: PropTypes.string,
  getAccounts: PropTypes.func.isRequired,
  newAccount: PropTypes.func.isRequired,
  newExternal: PropTypes.func.isRequired
};

const mapStateToProps = (state) => ({
  accounts: state.aws.accounts.all,
  external: state.aws.accounts.external
});

const mapDispatchToProps = (dispatch) => ({
  getAccounts: () => {
    dispatch(Actions.AWS.Accounts.getAccounts())
  },
  newAccount: (account) => {
    dispatch(Actions.AWS.Accounts.newAccount(account))
  },
  newExternal: () => {
    dispatch(Actions.AWS.Accounts.newExternal())
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(AccountsContainer);
