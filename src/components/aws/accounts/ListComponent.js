import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Misc from '../../misc';
import Form from './FormComponent';
import Bills from './bills';

const Panel = Misc.Panel;

export class ListItem extends Component {

  constructor(props) {
    super(props);
    this.state = {
      editForm: false
    };
    this.toggleEditForm = this.toggleEditForm.bind(this);
    this.editAccount = this.editAccount.bind(this);
    this.deleteAccount = this.deleteAccount.bind(this);
    this.bills = [{
      "bucket": "s3://te.st",
      "path": "/path/to/bills"
    },{
      "bucket": "s3://another.test",
      "path": "/another/path"
    }]
  }

  toggleEditForm = (e) => {
    e.preventDefault();
    const editForm = !this.state.editForm;
    this.setState({ editForm });
  };

  editAccount = (body) => {
    this.setState({ editForm: false });
    this.props.accountActions.edit(body);
  };

  deleteAccount = (e) => {
    e.preventDefault();
    this.props.accountActions.delete(this.props.account.id);
  };

  render() {

    const actions = (!this.state.editForm ? (
      <div className="btn-group btn-group-xs pull-right" role="group">

        <button type="button" className="btn btn-default edit" onClick={this.toggleEditForm}>
          <span className="glyphicon glyphicon-pencil" aria-hidden="true"/> Edit
        </button>

        <button type="button" className="btn btn-default delete" onClick={this.deleteAccount}>
          <span className="glyphicon glyphicon-remove" aria-hidden="true"/> Delete
        </button>

      </div>
    ) : (
      <div className="btn-group btn-group-xs pull-right" role="group">
        <button type="button" className="btn btn-default hide" onClick={this.toggleEditForm}>
          <span className="glyphicon glyphicon-remove" aria-hidden="true"/> Hide
        </button>
      </div>
    ));

    const editForm = (this.state.editForm ? (
      <div>

        <Form submit={this.editAccount} account={this.props.account} billActions={this.props.billActions}/>

        <Panel title="Bills locations" collapsible defaultCollapse >

          <Bills.List
            account={this.props.account.id}
            bills={this.bills}
            new={this.props.billActions.new}
            edit={this.props.billActions.edit}
            delete={this.props.billActions.delete}
          />

          <Bills.Form
            account={this.props.account.id}
            submit={this.props.billActions.new}
          />

        </Panel>

      </div>
    ) : null);

    return (
      <li className="account list-group-item">

        <span className="account-name">
          {this.props.account.pretty || this.props.account.roleArn}
        </span>

        {actions}

        {editForm}

      </li>
    );
  }

}

ListItem.propTypes = {
  account: PropTypes.shape({
    id: PropTypes.number.isRequired,
    roleArn: PropTypes.string.isRequired,
    pretty: PropTypes.string,
    bills: PropTypes.arrayOf(
      PropTypes.shape({
        bucket: PropTypes.string.isRequired,
        path: PropTypes.string.isRequired
      })
    ),
  }),
  accountActions: PropTypes.shape({
    edit: PropTypes.func.isRequired,
    delete: PropTypes.func.isRequired,
  }).isRequired,
  billActions: PropTypes.shape({
    new: PropTypes.func.isRequired,
    edit: PropTypes.func.isRequired,
    delete: PropTypes.func.isRequired,
  }).isRequired,
};

// List Component for AWS Accounts
class ListComponent extends Component {

  render() {
    let noAccounts = (!this.props.accounts || !this.props.accounts.length ? <div className="alert alert-warning" role="alert">No account available</div> : "");
    return (
      <ul className="accounts list-group list-group-flush">
        {noAccounts}
        {this.props.accounts.map((account) => (
          <ListItem
            key={account.id}
            account={account}
            accountActions={this.props.accountActions}
            billActions={this.props.billActions}
          />
        ))}
      </ul>
    );
  }

}

ListComponent.propTypes = {
  accounts: PropTypes.arrayOf(
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
  accountActions: PropTypes.shape({
    edit: PropTypes.func.isRequired,
    delete: PropTypes.func.isRequired,
  }).isRequired,
  billActions: PropTypes.shape({
    new: PropTypes.func.isRequired,
    edit: PropTypes.func.isRequired,
    delete: PropTypes.func.isRequired,
  }).isRequired,
};

export default ListComponent;
