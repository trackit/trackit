import React, { Component } from 'react';
import { connect } from 'react-redux';

import Components from '../../../components';
import Actions from "../../../actions";

const List = Components.AWS.Accounts.List;

// MainContainer Component
class AccountsContainer extends Component {

  componentWillMount() {
    this.props.getAccounts();
  }

  render() {
    console.log(this.props.accounts);
    return (
      <div>
        <List/>
        {this.props.accounts.map((item) => (JSON.stringify(item)))}
      </div>
    );
  }

}

const mapStateToProps = (state) => ({accounts: state.aws.accounts});

const mapDispatchToProps = (dispatch) => ({
  getAccounts: () => {
    dispatch(Actions.AWS.Accounts.getAccounts())
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(AccountsContainer);
