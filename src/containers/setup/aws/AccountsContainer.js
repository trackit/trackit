import React, { Component } from 'react';
import { connect } from 'react-redux';

import Components from '../../../components';
import Actions from "../../../actions";

const List = Components.AWS.Accounts.List;
const Form = Components.AWS.Accounts.Form;

// AccountsContainer Component
class AccountsContainer extends Component {

  constructor(props) {
    super(props);
    this.addAccount = this.addAccount.bind(this);
  }

  componentWillMount() {
    this.props.getAccounts();
  }

  addAccount = (account) => {
    this.props.newAccount(account);
  };

  render() {
    return (
      <div className="panel panel-default">
        <div className="panel-heading">
          <h3 className="panel-title">AWS Accounts</h3>
        </div>
        <div className="panel-body">
          <List accounts={this.props.accounts}/>
          <Form submit={this.addAccount}/>
        </div>
      </div>
    );
  }

}

const mapStateToProps = (state) => ({
  accounts: state.aws.accounts
});

const mapDispatchToProps = (dispatch) => ({
  getAccounts: () => {
    dispatch(Actions.AWS.Accounts.getAccounts())
  },
  newAccount: (account) => {
    dispatch(Actions.AWS.Accounts.newAccount(account))
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(AccountsContainer);
