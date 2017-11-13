import React, { Component } from 'react';
import PropTypes from 'prop-types';

import Form from './FormComponent';

class ListItem extends Component {

  constructor(props) {
    super(props);
    this.state = {
      editForm: false
    };
    this.showEditForm = this.showEditForm.bind(this);
    this.editAccount = this.editAccount.bind(this);
    this.deleteAccount = this.deleteAccount.bind(this);
  }

  showEditForm = (e) => {
    e.preventDefault();
    this.setState({ editForm: true });
  };

  editAccount = (body) => {
    console.log(body);
    this.setState({ editForm: false });
  };

  deleteAccount = (e) => {
    e.preventDefault();
    this.props.delete(this.props.account.id);
  };

  render() {

    const editButton = (!this.state.editForm ? (
      <button type="button" className="btn btn-default" onClick={this.showEditForm}>
        <span className="glyphicon glyphicon-pencil" aria-hidden="true"/> Edit
      </button>
    ) : null);

    const editForm = (this.state.editForm ? (
      <Form submit={this.editAccount} account={this.props.account}/>
    ) : null);

    return (
      <li className="account list-group-item">

        <span className="account-name">
          {this.props.account.pretty || this.props.account.roleArn}
        </span>

        <div className="btn-group btn-group-xs pull-right" role="group">

          {editButton}

          <button type="button" className="btn btn-default" onClick={this.deleteAccount}>
            <span className="glyphicon glyphicon-remove" aria-hidden="true"/> Delete
          </button>

        </div>

        {editForm}

      </li>
    );
  }

}

ListItem.propTypes = {
  account: PropTypes.shape({
    id: PropTypes.number.isRequired,
    roleArn: PropTypes.string.isRequired,
    userId: PropTypes.number.isRequired,
    pretty: PropTypes.string
  }),
  delete: PropTypes.func.isRequired
};

// List Component for AWS Accounts
class ListComponent extends Component {

  render() {
    let noAccounts = (!this.props.accounts.length ? <div className="alert alert-warning" role="alert">No account available</div> : "");
    return (
      <ul className="accounts list-group list-group-flush">
        {noAccounts}
        {this.props.accounts.map((account) => (<ListItem key={account.id} account={account} delete={this.props.delete}/>))}
      </ul>
    );
  }

}

ListComponent.propTypes = {
  accounts: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.number.isRequired,
      roleArn: PropTypes.string.isRequired,
      userId: PropTypes.number.isRequired,
      pretty: PropTypes.string
    })
  ),
  delete: PropTypes.func.isRequired
};

export default ListComponent;
