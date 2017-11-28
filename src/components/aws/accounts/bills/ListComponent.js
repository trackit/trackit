import React, { Component } from 'react';
import List, {
  ListItem,
  ListItemText,
} from 'material-ui/List';
import Misc from '../../../misc';
import PropTypes from 'prop-types';
import Form from './FormComponent';

const Dialog = Misc.Dialog;
const DeleteConfirmation = Misc.DeleteConfirmation;

class Item extends Component {

  constructor(props) {
    super(props);
    this.editBill = this.editBill.bind(this);
    this.deleteBill = this.deleteBill.bind(this);
  }

  editBill = (body) => {
    this.props.edit(body);
  };

  deleteBill = () => {
    this.props.delete(this.props.bill.id);
  };

  render() {

    return (
      <ListItem divider>

        <ListItemText
          disableTypography
          primary={this.props.bill.bucket + this.props.bill.path}
        />

        <div>

          <div className="inline-block">
            <Form
              account={this.props.account}
              bill={this.props.bill}
              submit={this.editBill}
            />
          </div>
          &nbsp;
          <div className="inline-block">
            <DeleteConfirmation entity="account" confirm={this.deleteBill}/>
          </div>

        </div>

      </ListItem>
    );
  }

}

Item.propTypes = {
  account: PropTypes.number.isRequired,
  bill: PropTypes.shape({
    bucket: PropTypes.string.isRequired,
    path: PropTypes.string.isRequired
  }),
  edit: PropTypes.func.isRequired,
  delete: PropTypes.func.isRequired
};

// List Component for AWS Accounts
class ListComponent extends Component {

  render() {
    let noBills = (!this.props.bills.length ? <div className="alert alert-warning" role="alert">No bills available</div> : "");
    return (
      <Dialog
        buttonName="Bills locations"
        title="Bills locations"
        secondActionName="Close"
      >

        <Form
          account={this.props.account}
          submit={this.editBill}
        />

        <List>
          {noBills}
          {this.props.bills.map((bill, index) => (
            <Item
              key={index}
              bill={bill}
              account={this.props.account}
              edit={this.props.edit}
              delete={this.props.delete}/>
          ))}
        </List>

      </Dialog>
    );
  }

}

ListComponent.propTypes = {
  account: PropTypes.number.isRequired,
  bills: PropTypes.arrayOf(
    PropTypes.shape({
      bucket: PropTypes.string.isRequired,
      path: PropTypes.string.isRequired
    })
  ),
  new: PropTypes.func.isRequired,
  edit: PropTypes.func.isRequired,
  delete: PropTypes.func.isRequired
};

export default ListComponent;
