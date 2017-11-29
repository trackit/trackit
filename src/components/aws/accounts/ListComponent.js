import React, { Component } from 'react';
import List, {
  ListItem,
  ListItemText,
} from 'material-ui/List';
import Collapse from 'material-ui/transitions/Collapse';
import PropTypes from 'prop-types';
import Misc from '../../misc';
import Form from './FormComponent';
import Bills from './bills';

const DeleteConfirmation = Misc.DeleteConfirmation;

class Item extends Component {

  constructor(props) {
    super(props);
    this.state = {
      showBills: false
    };
    this.toggleBills = this.toggleBills.bind(this);
    this.editAccount = this.editAccount.bind(this);
    this.deleteAccount = this.deleteAccount.bind(this);
    this.test = this.test.bind(this);
    this.bills = [{
      "bucket": "s3://te.st",
      "path": "/path/to/bills"
    },{
      "bucket": "s3://another.test",
      "path": "/another/path"
    }]
  }

  toggleBills = (e) => {
    e.preventDefault();
    const showBills = !this.state.showBills;
    this.setState({ showBills });
  };

  editAccount = (body) => {
    this.setState({ editForm: false });
    this.props.accountActions.edit(body);
  };

  deleteAccount = () => {
    this.props.accountActions.delete(this.props.account.id);
  };

  test = (e) => {
    e.preventDefault();
    console.log("hehe");
  };

  render() {
    return (
      <div>

        <ListItem divider>

          <ListItemText
            disableTypography
            primary={this.props.account.pretty || this.props.account.roleArn}
          />

          <div>

            <div className="inline-block">
              <Bills.List
                account={this.props.account.id}
                bills={this.bills}
                new={this.props.billActions.new}
                edit={this.props.billActions.edit}
                delete={this.props.billActions.delete}
              />
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

        <Collapse in={this.state.showBills} transitionDuration="auto" unmountOnExit>

          <Bills.List
            account={this.props.account.id}
            bills={this.bills}
            new={this.props.billActions.new}
            edit={this.props.billActions.edit}
            delete={this.props.billActions.delete}
          />

        </Collapse>

      </div>
    );
  }

}

Item.propTypes = {
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
    let accounts = (this.props.accounts && this.props.accounts.length ? (
      this.props.accounts.map((account, index) => (
        <Item
          key={index}
          account={account}
          accountActions={this.props.accountActions}
          billActions={this.props.billActions}
        />
      ))
    ) : null);
    return (
      <List disablePadding>
        {noAccounts}
        {accounts}
      </List>
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
